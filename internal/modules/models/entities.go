package models

import "time"

// User represents an authenticated platform user.
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Photo     string `json:"photo"`
	CompanyID int64  `json:"company_id"`
}

// Company represents an organization owned by a user.
type Company struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Website   string    `json:"website"`
	Owner     int64     `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

// CompanyMembership combines company data with the current user's membership role.
type CompanyMembership struct {
	Company
	Role      string `json:"role"`
	IsDefault bool   `json:"is_default"`
}

// CompanyMember represents a person in a company who can be added to teams.
type CompanyMember struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Team represents a company team used for culture and matching calculations.
type Team struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CompanyID   int64   `json:"company_id"`
}

// MCPToken is a one-time token for candidate or member quiz submissions.
type MCPToken struct {
	ID         int64
	EntityType string
	EntityID   int64
	Token      string
	ExpiresAt  time.Time
	Used       bool
}
