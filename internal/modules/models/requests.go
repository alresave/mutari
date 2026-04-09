package models

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Website string `json:"website" binding:"omitempty,url"`
}

type CreateCompanyMemberRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddTeamMemberRequest struct {
	CompanyMemberID int64 `json:"company_member_id" binding:"required"`
}

type CreateCandidateRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	PhotoURL string `json:"photo_url"`
}

type GenerateTokenRequest struct {
	CandidateID     *int64 `json:"candidate_id"`
	CompanyMemberID *int64 `json:"company_member_id"`
}

type SubmitAnswersRequest struct {
	Token   string        `json:"token" binding:"required"`
	Answers []AnswerInput `json:"answers" binding:"required"`
}

type SubmitMCPRequest struct {
	MemberID int64    `json:"member_id"`
	Answers  []Answer `json:"answers"`
}
