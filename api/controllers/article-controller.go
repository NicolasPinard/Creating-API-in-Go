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

// ArticleController struct
type ArticleController struct {
	articleRepo *repositories.ArticleRepository
}

// Init initialize the ArticleController
func (a *ArticleController) Init(db *sql.DB) {

	a.articleRepo = &repositories.ArticleRepository{}
	a.articleRepo.Init(db)
}

func (a *ArticleController) GetAllArticles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: GET all Articles endpoint")
	articleList, err := a.articleRepo.ScanArticles()
	if err != nil {
		log.Printf("Error while scaning all articles: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(articleList)
}

func (a *ArticleController) GetArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint hit: GET article endpoint with id " + key)
	keyID, err := strconv.ParseInt(key, 10, 64)
	if err != nil || keyID < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "the article id given is not a positive integer"}`))
		return
	}

	article, _ := a.articleRepo.QueryArticle(keyID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(article)
}

func (a *ArticleController) CreateArticle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: POST new article")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var article models.Article
	err := decoder.Decode(&article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	articleID, err := a.articleRepo.InsertArticle(&article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "unable to persist article"}`))
		return
	}
	fmt.Printf("Created article %d\n", articleID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success": "created article"}`))
}
