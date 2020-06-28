package models

// Experience is a work or academic experience
type Experience struct {
	ID               int64  `json:"ID"`
	Title            string `json:"Title"`
	Responsabilities string `json:"Responsabilities"`
	Achievements     string `json:"Achievements"`
}

// Experiences is an array of Experience
type Experiences []Experience
