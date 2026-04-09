package models

type MCPAnswer struct {
	QuestionID int     `json:"question_id"`
	Block      int     `json:"block"`
	Value      float64 `json:"value"`
}

type MCPQuestion struct {
	ID               int64  `json:"id"`
	Block            int64  `json:"block"`
	BlockName        string `json:"block_name"`
	BlockDescription string `json:"block_description"`
	Question         string `json:"question"`
	Scale            string `json:"scale"`
}

type Question struct {
	ID             int
	DimensionIndex int
	Reverse        bool
	IsControl      bool
}

type Answer struct {
	QuestionID int64   `json:"question_id"`
	Block      int     `json:"block"`
	Value      float64 `json:"value"`
}

type AnswerInput struct {
	QuestionID int64   `json:"question_id"`
	Value      float64 `json:"value"`
}
