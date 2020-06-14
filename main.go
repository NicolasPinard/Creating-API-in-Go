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

// Article forms the most basic datatype of this API
type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
	URL     string `json:"URL"`
}

// Articles is an array of Article
type Articles []Article

// Client wraps a pool of Sqlite connections.
type Client struct {
	*sql.DB
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: All Articles endpoint")
	dbClient := NewDBClient("./articles.db")
	defer dbClient.Close()
	articleList, err := dbClient.ScanArticles()
	if err != nil {
		log.Printf("Error while scaning all articles: %v\n", err)
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
	}

	dbClient := NewDBClient("./articles.db")
	defer dbClient.Close()
	article, _ := dbClient.QueryArticle(keyID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("Successfully created %d row(s).\n", rows)

	return execResult.LastInsertId()
}

// QueryArticle retrieves an article from the database
func (c *Client) QueryArticle(id int64) (*Article, error) {
	row := c.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='articles'")
	var tableName string
	err := row.Scan(&tableName)
	if err != nil {
		log.Fatalln("Unable to find table 'articles'")
		return nil, err
	}

	statement, _ := c.Prepare("SELECT * FROM articles WHERE id=$1")
	articleRow := statement.QueryRow(id)

	var articleID int64
	var title string
	var desc string
	var content string
	var url string
	err = articleRow.Scan(&articleID, &title, &desc, &content, &url)
	a := Article{Title: title, Desc: desc, Content: content, URL: url}
	if err != nil {
		log.Printf("Unable to find article with id %d\n", id)
		return nil, err
	}
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
