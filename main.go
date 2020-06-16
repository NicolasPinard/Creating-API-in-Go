package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

// Article is one of the page of my blog
type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
	URL     string `json:"URL"`
}

// Articles is an array of Article
type Articles []Article

// Experience is a work or academic experience
type Experience struct {
	id               int64  `json:"ID"`
	Title            string `json:"Title"`
	Responsabilities string `json:"Responsabilities"`
	Achievements     string `json:"Achievements"`
}

// Experiences is an array of Experience
type Experiences []Experience

// Client wraps a pool of Sqlite connections.
type Client struct {
	*sql.DB
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: GET all Articles endpoint")
	dbClient := NewDBClient("./articles.db")
	defer dbClient.Close()
	articleList, err := dbClient.ScanArticles()
	if err != nil {
		log.Printf("Error while scaning all articles: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articleList)
}

func getArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint hit: GET article endpoint with id " + key)
	keyID, err := strconv.ParseInt(key, 10, 64)
	if err != nil || keyID < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "the article id given is not a positive integer"}`))
		return
	}

	dbClient := NewDBClient("./articles.db")
	defer dbClient.Close()
	article, _ := dbClient.QueryArticle(keyID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func createArticle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: POST new article")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var article Article
	err := decoder.Decode(&article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	client := NewDBClient("./articles.db")
	defer client.Close()
	articleID, err := client.InsertArticle(&article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "unable to persist article"}`))
		return
	}
	fmt.Printf("Created article %d\n", articleID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success": "created article"}`))
}

func getAllExperiences(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: GET all experiences endpoint")
	dbClient := NewDBClient("./portfolio.db")
	defer dbClient.Close()
	experienceList, err := dbClient.ScanExperiences()
	if err != nil {
		log.Printf("Error while scaning all experiences: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(experienceList)
}

func getExperience(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint hit: GET experience endpoint with id " + key)
	keyID, err := strconv.ParseInt(key, 10, 64)
	if err != nil || keyID < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "the experience id given must be a positive integer"}`))
		return
	}

	dbClient := NewDBClient("./portfolio.db")
	defer dbClient.Close()
	experience, _ := dbClient.QueryExperience(keyID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(experience)
}

func createExperience(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: POST new experience")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var exp Experience
	err := decoder.Decode(&exp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	client := NewDBClient("./portfolio.db")
	defer client.Close()
	experienceID, err := client.InsertExperience(&exp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "unable to persist experience"}`))
		return
	}
	exp.id = experienceID
	fmt.Printf("Created experience %d\n", experienceID)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"success": "created experience"}`))
}

func homePage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to the Homepage!")
	fmt.Println("Endpoint hit: homepage")
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/articles", getAllArticles).Methods(http.MethodGet)
	r.HandleFunc("/articles", createArticle).Methods(http.MethodPost)
	r.HandleFunc("/articles/{id}", getArticle).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// NewDBClient creates a Sqlite3 Client
func NewDBClient(path string) *Client {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error while connecting to sqlite3 database: %v", err)
	}
	return &Client{db}
}

// InsertArticle persists the article passed to the database client
func (c *Client) InsertArticle(a *Article) (int64, error) {

	statement, err := c.Prepare("CREATE TABLE IF NOT EXISTS articles (id INTEGER PRIMARY KEY, title TEXT, desc TEXT, content TEXT, url TEXT)")
	if err != nil {
		log.Fatalf("Unable to create table: %v", err)
		return -1, err
	}
	statement.Exec()
	statement, _ = c.Prepare("INSERT INTO articles (title, desc, content, url) VALUES (?, ?, ?, ?)")
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
func (c *Client) QueryArticle(id int64) (*Article, error) {

	row := c.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='articles'")
	var tableName string
	err := row.Scan(&tableName)
	if err != nil {
		log.Println("Unable to find table 'articles'")
	}

	statement, _ := c.Prepare("SELECT * FROM articles WHERE id=$1")
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
	a := Article{Title: title, Desc: desc, Content: content, URL: url}
	return &a, nil
}

// ScanArticles returns all articles from the database
func (c *Client) ScanArticles() (Articles, error) {

	rows, err := c.Query("SELECT * FROM articles")
	if err != nil {
		return nil, err
	}

	var articleList Articles
	var articleID int64
	var title string
	var desc string
	var content string
	var url string
	for rows.Next() {
		rows.Scan(&articleID, &title, &desc, &content, &url)
		articleList = append(articleList, Article{Title: title, Desc: desc, Content: content, URL: url})
	}

	return articleList, nil
}

// InsertExperience persists the experience passed to the database client
func (c *Client) InsertExperience(e *Experience) (int64, error) {

	statement, err := c.Prepare("CREATE TABLE IF NOT EXISTS experience (id INTEGER PRIMARY KEY, title TEXT, responsabilities TEXT, achievements TEXT)")
	if err != nil {
		log.Fatalf("Unable to create table experience: %v", err)
	}
	statement.Exec()
	statement, _ = c.Prepare("INSERT INTO experience (title, responsabilities, achievements) VALUES (?, ?, ?)")
	execResult, err := statement.Exec(e.Title, e.Responsabilities, e.Achievements)
	if err != nil {
		log.Printf("Got an error while trying to insert experience: %v", err)
		return 0, err
	}
	rows, _ := execResult.RowsAffected()
	fmt.Printf("Successfully create %d experience(s).\n", rows)

	return execResult.LastInsertId()
}

// QueryExperience retrieves an experience from the database
func (c *Client) QueryExperience(e int64) (*Experience, error) {

	row := c.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='experience'")
	err := row.Scan()
	if err != nil {
		log.Println("Unable to find table experience")
	}

	statement, _ := c.Prepare("SELECT * FROM experience WHERE id=$1")
	experienceRow := statement.QueryRow(e)

	var experienceID int64
	var title string
	var responsabilities string
	var achievements string
	err = experienceRow.Scan(&experienceID, &title, &responsabilities, &achievements)
	if err != nil {
		log.Printf("Unable to find experience with id %d\n", e)
		return nil, err
	}
	exp := Experience{id: experienceID, Title: title, Responsabilities: responsabilities, Achievements: achievements}
	return &exp, nil
}

// ScanExperiences returns all experiences from the database
func (c *Client) ScanExperiences() (Experiences, error) {

	rows, err := c.Query("SELECT * FROM experience")
	if err != nil {
		return nil, err
	}

	var experienceList Experiences
	var experienceID int64
	var title string
	var responsabilities string
	var achievements string
	for rows.Next() {
		rows.Scan(&experienceID, &title, &responsabilities, &achievements)
		experienceList = append(experienceList, Experience{id: experienceID, Title: title, Responsabilities: responsabilities, Achievements: achievements})
	}

	return experienceList, nil
}
