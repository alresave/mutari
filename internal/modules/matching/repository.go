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
	UserID      int64  `json:"user_id"`
	Score       int    `json:"score"`
	TopBlock    string `json:"top_block"`
	Explanation string `json:"explanation"`
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

		for b := 0; b < 8; b++ {

			start := b * 8
			end := start + 8

			sim := blockSimilarity(vecA, vecB, start, end)

			if sim > topBlockScore {
				topBlockScore = sim
				topBlock = blockNames[b]
			}

			w := blockWeights[b]

			weightedScore += sim * w
			totalWeight += w
		}

		structural := weightedScore / totalWeight

		finalScore := structural * math.Sqrt(baseConf*conf)

		percentage := int(math.Round(finalScore * 100))

		explanation := explainMatch(topBlock, finalScore)

		candidates = append(candidates, Candidate{
			UserID:      id,
			Score:       percentage,
			TopBlock:    topBlock,
			Explanation: explanation,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	return candidates, nil
}
