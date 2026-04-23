package mcp

import (
	"context"
	"fmt"
	"math"
	"strings"

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
func (r *Repository) StoreCompanyMemberAnswers(ctx context.Context, memberID int64, answers []models.MCPAnswer) error {

	query := `
	INSERT INTO company_member_mcp_answers
	(company_member_id,question_id,value)
	VALUES ($1,$2,$3)
	ON CONFLICT (company_member_id, question_id)
	DO UPDATE SET value = EXCLUDED.value
	`

	for _, a := range answers {

		_, err := r.db.Exec(ctx, query,
			memberID,
			a.QuestionID,
			a.Value,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
func (r *Repository) BuildCompanyMemberProfile(ctx context.Context, memberID int64) error {

	query := `
	SELECT q.block, a.value
	FROM company_member_mcp_answers a
	JOIN mcp_questions q ON q.id = a.question_id
	WHERE a.company_member_id=$1
	`

	rows, err := r.db.Query(ctx, query, memberID)
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

		blockIndex := block - 1
		if blockIndex >= 0 && blockIndex < 8 {
			blockSum[blockIndex] += value
			blockCount[blockIndex]++
		}
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

		archetype := archetypes[int(memberID)%len(archetypes)]
		base := archetype[block]

		start := block * 8

		for i := 0; i < 8; i++ {
			gradient := (float64(i) - 3.5) * 0.30

			val := (base * 0.6) + (avg * 0.4) + gradient

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
	INSERT INTO company_member_mcp_profiles(company_member_id,effective_vector)
	VALUES($1,$2)
	ON CONFLICT(company_member_id)
	DO UPDATE SET effective_vector=EXCLUDED.effective_vector
	`

	_, err = r.db.Exec(ctx, queryInsert, memberID, pgVector)

	return err
}

func (r *Repository) HasMCPProfile(ctx context.Context, userID int64) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM company_member_mcp_profiles cmp
		JOIN company_members cm ON cm.id = cmp.company_member_id
		JOIN companies c ON c.id = cm.company_id
		WHERE c.owner=$1
	)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)

	return exists, err
}

func (r *Repository) GetAllCompanyMemberIDs(ctx context.Context) ([]int64, error) {

	rows, err := r.db.Query(ctx, `SELECT id FROM company_members`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
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

func cosineSimilarity(a, b []float32) float32 {

	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float32

	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (float32(math.Sqrt(float64(normA))) *
		float32(math.Sqrt(float64(normB))))
}
func cosineDistance(a, b []float32) float32 {

	if len(a) != len(b) {
		return 1 // distancia máxima
	}

	var dot float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {

		av := a[i]
		bv := b[i]

		dot += av * bv
		normA += av * av
		normB += bv * bv
	}

	if normA == 0 || normB == 0 {
		return 1
	}

	// 🔥 conversión correcta
	denominator := float32(
		math.Sqrt(float64(normA)) *
			math.Sqrt(float64(normB)),
	)

	return 1 - (dot / denominator)
}

func (r *Repository) SaveAnswers(
	ctx context.Context,
	candidateID int64,
	answers []models.Answer,
) error {

	if len(answers) == 0 {
		return nil
	}

	// 🔥 construir query dinámica
	query := `
		INSERT INTO candidate_mcp_answers (candidate_id, question_id, value)
		VALUES 
	`

	args := []interface{}{}
	values := []string{}

	for i, a := range answers {

		// ($1, $2, $3), ($4, $5, $6)...
		idx := i * 3

		values = append(values,
			fmt.Sprintf("($%d, $%d, $%d)",
				idx+1, idx+2, idx+3),
		)

		args = append(args,
			candidateID,
			a.QuestionID,
			a.Value,
		)
	}

	query += strings.Join(values, ",")
	query += `
		ON CONFLICT (candidate_id, question_id)
		DO UPDATE SET value = EXCLUDED.value
	`

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) SaveMemberAnswers(
	ctx context.Context,
	memberID int64,
	answers []models.AnswerInput,
) error {

	if len(answers) == 0 {
		return nil
	}

	query := `
		INSERT INTO company_member_mcp_answers (company_member_id, question_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (company_member_id, question_id)
		DO UPDATE SET value = EXCLUDED.value
	`

	for _, a := range answers {
		if _, err := r.db.Exec(ctx, query, memberID, a.QuestionID, a.Value); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) UpsertCandidateProfile(
	ctx context.Context,
	candidateID int64,
	vector []float32,
) error {
	pgVector := pgvector.NewVector(vector)

	query := `
		INSERT INTO candidate_mcp_profiles (candidate_id, effective_vector)
		VALUES ($1, $2)
		ON CONFLICT (candidate_id)
		DO UPDATE SET effective_vector = EXCLUDED.effective_vector
	`

	_, err := r.db.Exec(ctx, query, candidateID, pgVector)
	return err
}

func (r *Repository) UpsertCompanyMemberProfile(
	ctx context.Context,
	memberID int64,
	vector []float32,
) error {
	pgVector := pgvector.NewVector(vector)

	query := `
		INSERT INTO company_member_mcp_profiles (company_member_id, effective_vector)
		VALUES ($1, $2)
		ON CONFLICT (company_member_id)
		DO UPDATE SET effective_vector = EXCLUDED.effective_vector
	`

	_, err := r.db.Exec(ctx, query, memberID, pgVector)
	return err
}
