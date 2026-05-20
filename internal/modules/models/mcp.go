package models

// MCPAnswer stores a submitted answer with its resolved MCP block.
type MCPAnswer struct {
	QuestionID int     `json:"question_id"`
	Block      int     `json:"block"`
	Value      float64 `json:"value"`
}

// MCPQuestion represents a quiz question with its block metadata.
type MCPQuestion struct {
	ID               int64  `json:"id"`
	Block            int64  `json:"block"`
	BlockName        string `json:"block_name"`
	BlockDescription string `json:"block_description"`
	Question         string `json:"question"`
	Scale            string `json:"scale"`
}

// Question maps a quiz prompt to a vector dimension for identity encoding.
type Question struct {
	ID             int
	DimensionIndex int
	Reverse        bool
	IsControl      bool
}

// Answer stores a normalized quiz answer used by the matching pipeline.
type Answer struct {
	QuestionID int64   `json:"question_id"`
	Block      int     `json:"block"`
	Value      float64 `json:"value"`
}

// AnswerInput is the API payload shape for raw submitted quiz answers.
type AnswerInput struct {
	QuestionID int64   `json:"question_id"`
	Value      float64 `json:"value"`
}
