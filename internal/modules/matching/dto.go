package matching

type SubmitMCPRequest struct {
	UserID  int64    `json:"user_id"`
	Answers []Answer `json:"answers"`
}
type MCPQuestion struct {
	ID        int    `json:"id"`
	Block     int    `json:"block"`
	BlockName string `json:"block_name"`
	Question  string `json:"question"`
	Scale     string `json:"scale"`
}
