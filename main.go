package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/NicolasPinard/PortfolioWebsite/api/controllers"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

// Client wraps a pool of Sqlite connections.
type Client struct {
	*sql.DB
}

func homePage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to the Homepage!")
	fmt.Println("Endpoint hit: homepage")
}

func main() {

	db := InitDB("./articles.db")
	defer db.Close()
	articleController := controllers.ArticleController{}
	articleController.Init(db)
	experienceController := controllers.ExperienceController{}
	experienceController.Init(db)

	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/articles", articleController.GetAllArticles).Methods(http.MethodGet)
	r.HandleFunc("/articles", articleController.CreateArticle).Methods(http.MethodPost)
	r.HandleFunc("/articles/{id}", articleController.GetArticle).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// InitDB initializes the connection to the database
func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error while connecting to sqlite3 database: %v", err)
	}
	return db
}

// NewDBClient creates a Sqlite3 Client
func NewDBClient(path string) *Client {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error while connecting to sqlite3 database: %v", err)
	}
	return &Client{db}
}
