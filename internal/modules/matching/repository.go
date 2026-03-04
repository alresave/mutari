package matching

import (
	"context"
	"math"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type Candidate struct {
	UserID               int64  `json:"user_id"`
	Score                int    `json:"score"`
	TopBlock             string `json:"top_block"`
	SecondaryBlock       string `json:"secondary_alignment"`
	CompatibilityType    string `json:"compatibility_type"`
	Explanation          string `json:"explanation"`
	PotentialFriction    string `json:"potential_friction"`
	CompatibilitySummary string `json:"compatibility_summary"`
}

var blockWeights = []float64{
	2.0, // Core Values
	1.8, // Relational Architecture
	1.2, // Lifestyle Rhythm
	1.4, // Growth
	1.3, // Cognitive Style
	1.4, // Ambition
	1.6, // Emotional Stability
	1.5, // Intimacy
}

var blockNames = []string{
	"Core Values",
	"Relational Architecture",
	"Lifestyle Rhythm",
	"Growth",
	"Cognitive Style",
	"Ambition",
	"Emotional Stability",
	"Intimacy",
}

func blockSimilarity(a, b []float32, start, end int) float64 {

	var dot float64
	var normA float64
	var normB float64

	for i := start; i < end; i++ {

		av := float64(a[i])
		bv := float64(b[i])

		dot += av * bv
		normA += av * av
		normB += bv * bv
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func explainMatch(block string, score float64) string {

	if score > 0.85 {
		return "Exceptional alignment in " + block
	}

	if score > 0.75 {
		return "Strong compatibility in " + block
	}

	if score > 0.65 {
		return "Moderate alignment in " + block
	}

	return "Some differences in " + block
}

func compatibilityType(block string) string {

	switch block {

	case "Core Values":
		return "Values Match"

	case "Relational Architecture":
		return "Relationship Structure Match"

	case "Lifestyle Rhythm":
		return "Lifestyle Match"

	case "Growth":
		return "Growth Alignment"

	case "Cognitive Style":
		return "Thinking Style Match"

	case "Ambition":
		return "Ambition Alignment"

	case "Emotional Stability":
		return "Emotional Compatibility"

	case "Intimacy":
		return "Bonding Compatibility"

	default:
		return "General Compatibility"
	}
}

func frictionExplanation(block string) string {

	switch block {

	case "Core Values":
		return "Possible differences in core life principles"

	case "Relational Architecture":
		return "Different expectations about relationship roles"

	case "Lifestyle Rhythm":
		return "Potential mismatch in daily routines or pace of life"

	case "Growth":
		return "Different priorities in personal development"

	case "Cognitive Style":
		return "Different ways of thinking or processing decisions"

	case "Ambition":
		return "Possible differences in long-term goals or drive"

	case "Emotional Stability":
		return "Different emotional regulation styles"

	case "Intimacy":
		return "Different bonding or closeness needs"

	default:
		return "Potential area of difference"
	}
}

func compatibilitySummary(primary, secondary, friction string) string {

	return "Strong alignment in " +
		primary +
		" with additional compatibility in " +
		secondary +
		". However, there may be differences in " +
		friction +
		"."
}

func (r *Repository) GetMCPMatches(ctx context.Context, userID int64) ([]Candidate, error) {

	query := `
	SELECT user_id, effective_vector, confidence_final
	FROM user_mcp_profiles
	WHERE user_id != $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var baseVec pgvector.Vector
	var baseConf float64

	err = r.db.QueryRow(ctx,
		`SELECT effective_vector, confidence_final
		 FROM user_mcp_profiles
		 WHERE user_id=$1`,
		userID).Scan(&baseVec, &baseConf)

	if err != nil {
		return nil, err
	}

	var candidates []Candidate

	for rows.Next() {

		var id int64
		var vec pgvector.Vector
		var conf float64

		if err := rows.Scan(&id, &vec, &conf); err != nil {
			return nil, err
		}

		vecA := baseVec.Slice()
		vecB := vec.Slice()

		var weightedScore float64
		var totalWeight float64

		var topBlockScore float64
		var topBlock string

		var secondScore float64
		var secondBlock string

		var lowestScore float64 = 1
		var frictionBlock string

		for b := 0; b < 8; b++ {

			start := b * 8
			end := start + 8

			sim := blockSimilarity(vecA, vecB, start, end)

			if sim > topBlockScore {

				secondScore = topBlockScore
				secondBlock = topBlock

				topBlockScore = sim
				topBlock = blockNames[b]

			} else if sim > secondScore {

				secondScore = sim
				secondBlock = blockNames[b]

			}

			if sim < lowestScore {
				lowestScore = sim
				frictionBlock = blockNames[b]
			}

			w := blockWeights[b]

			weightedScore += sim * w
			totalWeight += w
		}

		structural := weightedScore / totalWeight

		finalScore := structural * math.Sqrt(baseConf*conf)

		percentage := int(math.Round(finalScore * 100))

		explanation := explainMatch(topBlock, finalScore)

		ctype := compatibilityType(topBlock)

		friction := frictionExplanation(frictionBlock)

		summary := compatibilitySummary(topBlock, secondBlock, frictionBlock)

		candidates = append(candidates, Candidate{
			UserID:               id,
			Score:                percentage,
			TopBlock:             topBlock,
			SecondaryBlock:       secondBlock,
			CompatibilityType:    ctype,
			Explanation:          explanation,
			PotentialFriction:    friction,
			CompatibilitySummary: summary,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	return candidates, nil
}
