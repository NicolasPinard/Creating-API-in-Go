package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Article forms the most basic datatype of this API
type Article struct {
	Title   string `json:"Title`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
}

// Articles is an array of Article
type Articles []Article

var articles = Articles{
	Article{Title: "Test Title", Desc: "Lorem ipsum", Content: "Hello World!"},
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: All Articles endpoint")
	json.NewEncoder(w).Encode(articles)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	articles = append(articles, Article{Title: r.FormValue("Title"), Desc: r.FormValue("Desc"), Content: r.FormValue("Content")})
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success": "created"`))
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
	log.Fatal(http.ListenAndServe(":8080", r))
}
