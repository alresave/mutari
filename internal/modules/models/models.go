package models

import "time"

var RoleWeights = map[string][]float64{

	"engineer": {
		0.20, // Execution Drive
		0.30, // Cognitive Processing
		0.10, // Professional Values
		0.10, // Emotional Resilience
		0.20, // Learning Velocity
		0.05, // Collaboration
		0.03, // Operational Rhythm
		0.02, // Team Positioning
	},

	"manager": {
		0.10,
		0.05,
		0.15,
		0.15,
		0.03,
		0.25,
		0.02,
		0.25,
	},

	"startup": {
		0.25,
		0.15,
		0.02,
		0.10,
		0.25,
		0.05,
		0.15,
		0.03,
	},
}

var BlockNames = []string{
	"Execution Drive",
	"Cognitive Processing",
	"Professional Values",
	"Emotional Resilience",
	"Learning Velocity",
	"Collaboration Depth",
	"Operational Rhythm",
	"Team Positioning",
}
var BlockDifferences = map[string]string{
	"Execution Drive":      "Different levels of initiative, ownership, and focus on achieving results",
	"Cognitive Processing": "Different ways of thinking, analyzing information, and making decisions",
	"Professional Values":  "Possible differences in core professional values and priorities",
	"Emotional Resilience": "Different ways of handling stress, pressure, and emotional responses",
	"Learning Velocity":    "Differences in speed of learning, adaptability, and applying new knowledge",
	"Collaboration Depth":  "Different approaches to teamwork, communication, and shared problem-solving",
	"Operational Rhythm":   "Potential mismatch in work pace, structure, and execution style",
	"Team Positioning":     "Different preferences in roles and contribution within team dynamics",
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
type SubmitMCPRequest struct {
	UserID  int64    `json:"user_id"`
	Answers []Answer `json:"answers"`
}
type MCPQuestion struct {
	ID               int    `json:"id"`
	Block            int    `json:"block"`
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
	QuestionID int
	Value      float64
}
type Company struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Website   string    `json:"website"`
	Owner     int64     `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Website string `json:"website" binding:"required,url"`
}
type CreateCompanyMemberRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}
type CompanyMember struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
type AddTeamMemberRequest struct {
	CompanyMemberID int64 `json:"company_member_id" binding:"required"`
}

type CreateCandidateRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
