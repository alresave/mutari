package matching
type SubmitMCPRequest struct {
	UserID  int64    `json:"user_id"`
	Answers []Answer `json:"answers"`
}

