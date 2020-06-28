package models

// Article is one of the page of my blog
type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
	URL     string `json:"URL"`
}

// Articles is an array of Article
type Articles []Article
