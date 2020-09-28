package models

// Experience is a work or academic experience
type Experience struct {
	ID               int64  `json:"ID"`
	Title            string `json:"Title"`
	Responsibilities string `json:"Responsibilities"`
	Achievements     string `json:"Achievements"`
	Date			 string `json:"Date"`
}

// Experiences is an array of Experience
type Experiences []Experience
