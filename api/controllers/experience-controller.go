package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NicolasPinard/PortfolioWebsite/api/models"
	"github.com/NicolasPinard/PortfolioWebsite/api/repositories"
	"github.com/gorilla/mux"
)

// ExperienceController struct with functions for Experience struct
type ExperienceController struct {
	experienceRepo *repositories.ExperienceRepository
}

// Init initializes controller for ExperienceController struct
func (e *ExperienceController) Init(db *sql.DB) {

	e.experienceRepo = &repositories.ExperienceRepository{}
	e.experienceRepo.Init(db)
}

func (e *ExperienceController) GetAllExperiences(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: GET all experiences endpoint")
	experienceList, err := e.experienceRepo.ScanExperiences()
	if err != nil {
		log.Printf("Error while scaning all experiences: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(experienceList)
}

func (e *ExperienceController) GetExperience(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint hit: GET experience endpoint with id " + key)
	keyID, err := strconv.ParseInt(key, 10, 64)
	if err != nil || keyID < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "the experience id given must be a positive integer"}`))
		return
	}

	experience, _ := e.experienceRepo.QueryExperience(keyID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(experience)
}

func (e *ExperienceController) CreateExperience(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: POST new experience")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var exp models.Experience
	err := decoder.Decode(&exp)
	if err != nil {
		fmt.Printf("Got error %v while trying to decode object\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	experienceID, err := e.experienceRepo.InsertExperience(&exp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "unable to persist experience"}`))
		return
	}
	exp.ID = experienceID
	fmt.Printf("Created experience %d\n", experienceID)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"success": "created experience"}`))
}
