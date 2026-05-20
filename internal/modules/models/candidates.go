package models

import "time"

// Candidate2 is the dashboard-facing candidate record.
type Candidate2 struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PhotoURL  string    `json:"photo_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CandidateDB is the internal candidate row shape used by repository queries.
type CandidateDB struct {
	ID        int64
	CompanyID int64
	Name      string
	Email     string
	PhotoURL  string
	Status    string
	CreatedAt time.Time
}

// Candidatescore is a candidate record with an optional match score.
type Candidatescore struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PhotoURL  string    `json:"photo_url"`
	Status    string    `json:"status"`
	Score     float64   `json:"score,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
