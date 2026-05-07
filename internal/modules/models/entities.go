package models

import "time"

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Photo     string `json:"photo"`
	CompanyID int64  `json:"company_id"`
}

type Company struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Website   string    `json:"website"`
	Owner     int64     `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

type CompanyMembership struct {
	Company
	Role      string `json:"role"`
	IsDefault bool   `json:"is_default"`
}

type CompanyMember struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CompanyID   int64   `json:"company_id"`
}

type MCPToken struct {
	ID         int64
	EntityType string
	EntityID   int64
	Token      string
	ExpiresAt  time.Time
	Used       bool
}
