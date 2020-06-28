package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/NicolasPinard/PortfolioWebsite/api/models"
)

// ArticleRepository persistance struct for articles
type ArticleRepository struct {
	db *sql.DB
}

// Init initialize the ArticleRepository
func (repo *ArticleRepository) Init(db *sql.DB) {
	repo.db = db

	statement, err := repo.db.Prepare("CREATE TABLE IF NOT EXISTS articles (id INTEGER PRIMARY KEY, title TEXT, desc TEXT, content TEXT, url TEXT)")
	if err != nil {
		log.Fatalf("Unable to create table: %v", err)
	}
	statement.Exec()
}

// InsertArticle persists the article passed to the database client
func (repo *ArticleRepository) InsertArticle(a *models.Article) (int64, error) {

	statement, _ := repo.db.Prepare("INSERT INTO articles (title, desc, content, url) VALUES (?, ?, ?, ?)")
	execResult, err := statement.Exec(a.Title, a.Desc, a.Content, a.URL)
	if err != nil {
		return -1, err
	}
	rows, _ := execResult.RowsAffected()
	defer statement.Close()

	fmt.Printf("Successfully created %d article(s).\n", rows)

	return execResult.LastInsertId()
}

// QueryArticle retrieves an article from the database
func (repo *ArticleRepository) QueryArticle(id int64) (*models.Article, error) {

	row := repo.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='articles'")
	var tableName string
	err := row.Scan(&tableName)
	if err != nil {
		log.Println("Unable to find table 'articles'")
	}

	statement, _ := repo.db.Prepare("SELECT * FROM articles WHERE id=$1")
	articleRow := statement.QueryRow(id)

	var articleID int64
	var title string
	var desc string
	var content string
	var url string
	err = articleRow.Scan(&articleID, &title, &desc, &content, &url)
	if err != nil {
		log.Printf("Unable to find article with id %d\n", id)
		return nil, err
	}
	a := models.Article{Title: title, Desc: desc, Content: content, URL: url}
	return &a, nil
}

// ScanArticles returns all articles from the database
func (repo *ArticleRepository) ScanArticles() (models.Articles, error) {

	rows, err := repo.db.Query("SELECT * FROM articles")
	if err != nil {
		return nil, err
	}

	var articleList models.Articles
	var articleID int64
	var title string
	var desc string
	var content string
	var url string
	for rows.Next() {
		rows.Scan(&articleID, &title, &desc, &content, &url)
		articleList = append(articleList, models.Article{Title: title, Desc: desc, Content: content, URL: url})
	}

	return articleList, nil
}
