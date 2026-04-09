package models

type Candidate struct {
	MemberID int64 `json:"member_id"`

	Name  string `json:"name"`
	Photo string `json:"photo"`

	Score int `json:"score"`

	TopBlock       string `json:"top_block"`
	SecondaryBlock string `json:"secondary_alignment"`

	CompatibilityType string `json:"compatibility_type"`
	Explanation       string `json:"explanation"`
	PotentialFriction string `json:"potential_friction"`

	CompatibilitySummary  string `json:"compatibility_summary"`
	RelationshipStability string `json:"relationship_stability"`

	BlockScores map[string]int `json:"block_scores"`
	Role        string         `json:"role"`
}

type TeamFitSimulation struct {
	TeamID      int64          `json:"team_id"`
	CandidateID int64          `json:"candidate_id"`
	Score       int            `json:"score"`
	Role        string         `json:"suggested_role"`
	BlockScores map[string]int `json:"block_scores"`
}

type TeamVector struct {
	Vector    []float32 `json:"vector"`
	TeamScore float64   `json:"team_score"`
	Archetype string    `json:"archetype"`
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

type QuizResult struct {
	CandidateID int64
	BlockScores map[string]float64
}
