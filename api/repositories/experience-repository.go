package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/NicolasPinard/PortfolioWebsite/api/models"
)

// ExperienceRepository persistance struct for experiences
type ExperienceRepository struct {
	db *sql.DB
}

// Init initialize the ExperienceRepository
func (repo *ExperienceRepository) Init(db *sql.DB) {
	repo.db = db

	statement, err := repo.db.Prepare("CREATE TABLE IF NOT EXISTS experience (id INTEGER PRIMARY KEY, title TEXT, responsibilities TEXT, achievements TEXT, date TEXT)")
	if err != nil {
		log.Fatalf("Unable to create table experience: %v", err)
	}
	statement.Exec()
}

// InsertExperience persists the experience passed to the database client
func (repo *ExperienceRepository) InsertExperience(e *models.Experience) (int64, error) {

	statement, _ := repo.db.Prepare("INSERT INTO experience (title, responsibilities, achievements, date) VALUES (?, ?, ?, ?)")
	execResult, err := statement.Exec(e.Title, e.Responsibilities, e.Achievements, e.Date)
	if err != nil {
		log.Printf("Got an error while trying to insert experience: %v", err)
		return 0, err
	}
	rows, _ := execResult.RowsAffected()
	fmt.Printf("Successfully create %d experience(s).\n", rows)

	return execResult.LastInsertId()
}

// QueryExperience retrieves an experience from the database
func (repo *ExperienceRepository) QueryExperience(e int64) (*models.Experience, error) {

	row := repo.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='experience'")
	err := row.Scan()
	if err != nil {
		log.Println("Unable to find table experience")
	}

	statement, _ := repo.db.Prepare("SELECT * FROM experience WHERE id=$1")
	experienceRow := statement.QueryRow(e)

	var experienceID int64
	var title string
	var responsibilities string
	var achievements string
	var date string
	err = experienceRow.Scan(&experienceID, &title, &responsibilities, &achievements, &date)
	if err != nil {
		log.Printf("Unable to find experience with id %d\n", e)
		return nil, err
	}
	exp := models.Experience{ID: experienceID, Title: title, Responsibilities: responsibilities, Achievements: achievements, Date: date}
	return &exp, nil
}

// ScanExperiences returns all experiences from the database
func (repo *ExperienceRepository) ScanExperiences() (models.Experiences, error) {

	rows, err := repo.db.Query("SELECT * FROM experience")
	if err != nil {
		return nil, err
	}

	var experienceList models.Experiences
	var experienceID int64
	var title string
	var responsibilities string
	var achievements string
	var date string
	for rows.Next() {
		rows.Scan(&experienceID, &title, &responsibilities, &achievements, &date)
		experienceList = append(experienceList, models.Experience{ID: experienceID, Title: title, Responsibilities: responsibilities, Achievements: achievements, Date: date})
	}

	return experienceList, nil
}
