package mcp

import (
	"context"
	"math/rand/v2"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/pgvector/pgvector-go"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMCPQuestions(ctx context.Context) ([]models.MCPQuestion, error) {

	query := `
	SELECT
	q.id,
	q.block,
	b.name,
	b.description,
	q.question_text,
	q.scale
	FROM mcp_questions q
	JOIN mcp_blocks b ON q.block=b.id
	ORDER BY q.block,q.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.MCPQuestion

	for rows.Next() {

		var q models.MCPQuestion

		err := rows.Scan(
			&q.ID,
			&q.Block,
			&q.BlockName,
			&q.BlockDescription,
			&q.Question,
			&q.Scale,
		)

		if err != nil {
			return nil, err
		}

		questions = append(questions, q)
	}

	return questions, nil
}
func (r *Repository) StoreMCPAnswers(ctx context.Context, userID int64, answers []models.MCPAnswer) error {

	query := `
	INSERT INTO user_mcp_answers
	(user_id,question_id,block,value)
	VALUES ($1,$2,$3,$4)
	`

	for _, a := range answers {

		_, err := r.db.Exec(ctx, query,
			userID,
			a.QuestionID,
			a.Block,
			a.Value,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
func (r *Repository) BuildIdentityVector(ctx context.Context, userID int64) error {

	query := `
	SELECT block,value
	FROM user_mcp_answers
	WHERE user_id=$1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	blockSum := make(map[int]float64)
	blockCount := make(map[int]int)

	for rows.Next() {

		var block int
		var value float64

		if err := rows.Scan(&block, &value); err != nil {
			return err
		}

		blockSum[block] += value
		blockCount[block]++
	}

	// MCP archetype bases
	archetypes := [][]float64{

		// Leader
		{0.9, 0.7, 0.7, 0.7, 0.95, 0.8, 0.6, 0.8},

		// Strategist
		{0.8, 0.7, 0.8, 0.7, 0.7, 0.9, 0.7, 0.95},

		// Builder
		{0.8, 0.75, 0.7, 0.75, 0.9, 0.85, 0.65, 0.7},

		// Connector
		{0.85, 0.8, 0.9, 0.9, 0.65, 0.75, 0.9, 0.7},

		// Stabilizer
		{0.9, 0.85, 0.95, 0.85, 0.6, 0.7, 0.8, 0.65},
	}

	vector := make([]float32, 64)

	for block := 0; block < 8; block++ {

		avg := 0.5

		if blockCount[block] > 0 {
			avg = blockSum[block] / float64(blockCount[block])
		}

		archetype := archetypes[int(userID)%len(archetypes)]
		base := archetype[block]

		start := block * 8

		for i := 0; i < 8; i++ {
			noise := (rand.Float64() - 0.5) * 0.12
			gradient := (float64(i) - 3.5) * 0.30

			val := (base*0.6 + avg*0.4) + gradient + noise

			if val < 0 {
				val = 0
			}

			if val > 1 {
				val = 1
			}

			vector[start+i] = float32(val)
		}
	}

	pgVector := pgvector.NewVector(vector)

	queryInsert := `
	INSERT INTO user_mcp_profiles(user_id,effective_vector)
	VALUES($1,$2)
	ON CONFLICT(user_id)
	DO UPDATE SET effective_vector=EXCLUDED.effective_vector
	`

	_, err = r.db.Exec(ctx, queryInsert, userID, pgVector)

	return err
}

func (r *Repository) HasMCPProfile(ctx context.Context, userID int64) (bool, error) {

	query := `
	SELECT EXISTS(
	SELECT 1 FROMcompany_member_mcp_profiles 
	WHERE company_member_id=$1
	)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)

	return exists, err
}

func (r *Repository) GetTeamVectorRaw(
	ctx context.Context,
	teamID int64,
) ([]float64, error) {

	query := `
	SELECT effective_vector
	FROM team_profiles
	WHERE team_id = $1
	`

	var vec pgvector.Vector

	err := r.db.QueryRow(ctx, query, teamID).Scan(&vec)
	if err != nil {
		return nil, err
	}

	// 🔹 convertir []float32 → []float64
	raw := vec.Slice()

	vector := make([]float64, len(raw))

	for i, v := range raw {
		vector[i] = float64(v)
	}

	return vector, nil
}
func (r *Repository) GetCandidateVector(
	ctx context.Context,
	candidateID int64,
) ([]float32, error) {

	query := `
	SELECT effective_vector
	FROM candidate_mcp_profiles
	WHERE candidate_id = $1
	`

	var vec pgvector.Vector

	err := r.db.QueryRow(ctx, query, candidateID).Scan(&vec)
	if err != nil {
		return nil, err
	}

	return vec.Slice(), nil
}
