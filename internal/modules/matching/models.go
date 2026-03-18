package matching

var blockNames = []string{
	"Drive & Ownership",
	"Cognitive Style",
	"Work Values",
	"Emotional Regulation",
	"Learning Agility",
	"Collaboration",
	"Work Rythm ",
	"Team Architecture",
}

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Photo     string `json:"photo"`
	CompanyID int64  `json:"company_id"`
}

type TeamFitSimulation struct {
	TeamID      int64          `json:"team_id"`
	CandidateID int64          `json:"candidate_id"`
	Score       int            `json:"score"`
	Role        string         `json:"suggested_role"`
	BlockScores map[string]int `json:"block_scores"`
}
type Team struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CompanyID   int64   `json:"company_id"`
}
type TeamVector struct {
	Vector    []float32 `json:"vector"`
	TeamScore float64   `json:"team_score"`
	Archetype string    `json:"archetype"`
}
type Candidate struct {
	UserID int64 `json:"user_id"`

	Name  string `json:"name"`
	Photo string `json:"photo"`

	Score int `json:"score"`

	TopBlock       string `json:"top_block"`
	SecondaryBlock string `json:"secondary_alignment"`

	CompatibilityType string `json:"compatibility_type"`

	Explanation string `json:"explanation"`

	PotentialFriction string `json:"potential_friction"`

	CompatibilitySummary string `json:"compatibility_summary"`

	RelationshipStability string `json:"relationship_stability"`

	BlockScores map[string]int `json:"block_scores"`

	Role string `json:"role"`
}

type TeamCulture struct {
	TeamID      int64          `json:"team_id"`
	Archetype   string         `json:"archetype"`
	BlockScores map[string]int `json:"block_scores"`
}
type HiringProfile struct {
	TeamID            int64              `json:"team_id"`
	IdealProfile      map[string]float64 `json:"ideal_profile"`
	SuggestedRole     string             `json:"suggested_role"`
	RecommendedTraits []string           `json:"recommended_traits"`
}
type TeamGaps struct {
	TeamID             int64    `json:"team_id"`
	WeakBlocks         []string `json:"weak_blocks"`
	RecommendedProfile string   `json:"recommended_profile"`
}
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MCPAnswer struct {
	QuestionID int     `json:"question_id"`
	Block      int     `json:"block"`
	Value      float64 `json:"value"`
}
type SubmitAnswersRequest struct {
	UserID  int64       `json:"user_id"`
	Answers []MCPAnswer `json:"answers"`
}
