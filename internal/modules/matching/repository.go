package matching

import (
	"context"
	"errors"
	"math"

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

var blockNames = []string{
	"Core Values",
	"Lifestyle Rhythm",
	"Emotional Stability",
	"Relational Architecture",
	"Ambition",
	"Growth",
	"Intimacy",
	"Cognitive Style",
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

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func StabilityScore(score float64) string {

	switch {
	case score >= 0.9:
		return "Exceptional"
	case score >= 0.8:
		return "Strong"
	case score >= 0.6:
		return "Moderate"
	default:
		return "Uncertain"
	}
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
		return "Potential mismatch in work pace"

	case "Growth":
		return "Different priorities in personal development"

	case "Cognitive Style":
		return "Different ways of thinking"

	case "Ambition":
		return "Possible differences in long-term goals"

	case "Emotional Stability":
		return "Different emotional regulation styles"

	case "Intimacy":
		return "Different bonding needs"

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
	SELECT
	u.id,
	u.name,
	u.photo_url,
	p.effective_vector
	FROM user_mcp_profiles p
	JOIN users u ON u.id = p.user_id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vectors := make(map[int64][]float32)
	names := make(map[int64]string)
	photos := make(map[int64]string)

	for rows.Next() {

		var id int64
		var name string
		var photo string
		var vec pgvector.Vector

		if err := rows.Scan(&id, &name, &photo, &vec); err != nil {
			return nil, err
		}

		vectors[id] = vec.Slice()
		names[id] = name
		photos[id] = photo
	}

	source, ok := vectors[userID]
	if !ok {
		return nil, errors.New("user vector not found")
	}

	var results []Candidate

	for id, vec := range vectors {

		if id == userID {
			continue
		}

		blockScores := make(map[string]int)
		blockSim := make([]float64, 8)

		best := 0
		second := 1
		worst := 0

		for block := 0; block < 8; block++ {

			start := block * 8
			end := start + 8

			score := blockSimilarity(source, vec, start, end)

			blockSim[block] = score
			blockScores[blockNames[block]] = int(score * 100)
		}

		for i := 0; i < len(blockSim); i++ {

			if blockSim[i] > blockSim[best] {
				second = best
				best = i
			}

			if blockSim[i] < blockSim[worst] {
				worst = i
			}
		}

		score := (blockSim[best] + blockSim[second]) / 2

		c := Candidate{

			UserID: id,
			Name:   names[id],
			Photo:  photos[id],

			Score: int(score * 100),

			TopBlock:       blockNames[best],
			SecondaryBlock: blockNames[second],

			CompatibilityType: compatibilityType(blockNames[best]),

			Explanation: explainMatch(blockNames[best], score),

			PotentialFriction: frictionExplanation(blockNames[worst]),

			CompatibilitySummary: compatibilitySummary(
				blockNames[best],
				blockNames[second],
				blockNames[worst],
			),

			RelationshipStability: StabilityScore(score),

			BlockScores: blockScores,
		}

		results = append(results, c)
	}

	return results, nil
}

func (r *Repository) CreateUser(ctx context.Context, name, email, passwordHash string) (int64, error) {

	query := `
	INSERT INTO users (name,email,password_hash)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	var id int64

	err := r.db.QueryRow(ctx, query, name, email, passwordHash).Scan(&id)

	return id, err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (int64, string, string, error) {

	query := `
	SELECT id,password_hash,name
	FROM users
	WHERE email=$1
	`

	var id int64
	var hash string
	var name string

	err := r.db.QueryRow(ctx, query, email).Scan(&id, &hash, &name)

	return id, hash, name, err
}

func (r *Repository) HasMCPProfile(ctx context.Context, userID int64) (bool, error) {

	query := `
	SELECT EXISTS(
	SELECT 1 FROM user_mcp_profiles
	WHERE user_id=$1
	)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)

	return exists, err
}

func (r *Repository) GetMCPQuestions(ctx context.Context) ([]MCPQuestion, error) {

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

	var questions []MCPQuestion

	for rows.Next() {

		var q MCPQuestion

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

func (r *Repository) StoreMCPAnswers(ctx context.Context, userID int64, answers []MCPAnswer) error {

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

	vector := make([]float32, 64)

	for block := 0; block < 8; block++ {

		avg := 0.5

		if blockCount[block] > 0 {
			avg = blockSum[block] / float64(blockCount[block])
		}

		start := block * 8

		// cuánto se separa el bloque del centro psicológico
		spread := (avg - 0.5) * 0.8

		// pequeña variación por bloque
		blockBias := float64(block) * 0.01

		for i := 0; i < 8; i++ {

			// gradiente interno del bloque
			gradient := (float64(i) - 3.5) * 0.18

			val := 0.5 + spread + gradient + blockBias

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
func (r *Repository) GetAllUserIDs(ctx context.Context) ([]int64, error) {

	query := `SELECT id FROM users`

	rows, err := r.db.Query(ctx, query)
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

func (r *Repository) GetUserByID(ctx context.Context, userID int64) (int64, string, string, error) {

	query := `
	SELECT id, name, photo_url
	FROM users
	WHERE id=$1
	`

	var id int64
	var name string
	var photo string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&id,
		&name,
		&photo,
	)

	return id, name, photo, err
}

func (r *Repository) GetTeamMatches(ctx context.Context, teamID int64) ([]Candidate, error) {

	query := `
	SELECT
	u.id,
	u.name,
	u.photo_url,
	p.effective_vector,
	t.effective_vector
	FROM team_profiles t
	JOIN user_mcp_profiles p ON TRUE
	JOIN users u ON u.id = p.user_id
	WHERE t.team_id=$1
	AND u.id NOT IN (
		SELECT user_id FROM team_members WHERE team_id=$1
	)
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Candidate

	for rows.Next() {

		var id int64
		var name string
		var photo string
		var vec pgvector.Vector
		var teamVec pgvector.Vector

		err := rows.Scan(&id, &name, &photo, &vec, &teamVec)
		if err != nil {
			return nil, err
		}

		source := vec.Slice()
		target := teamVec.Slice()

		blockSim := make([]float64, 8)
		blockScores := make(map[string]int)

		best := 0
		second := 1
		worst := 0

		for block := 0; block < 8; block++ {

			start := block * 8
			end := start + 8

			score := blockSimilarity(source, target, start, end)

			blockSim[block] = score
			blockScores[blockNames[block]] = int(score * 100)
		}

		for i := 0; i < 8; i++ {

			if blockSim[i] > blockSim[best] {
				second = best
				best = i
			}

			if blockSim[i] < blockSim[worst] {
				worst = i
			}
		}

		score := (blockSim[best] + blockSim[second]) / 2

		stabilityWeights := []float64{
			0.25,
			0.10,
			0.20,
			0.20,
			0.05,
			0.05,
			0.10,
			0.05,
		}

		stability := 0.0

		for i := 0; i < 8; i++ {
			stability += blockSim[i] * stabilityWeights[i]
		}

		stabilityLabel := StabilityScore(stability)

		role := detectRole(blockScores)

		c := Candidate{

			UserID: id,
			Name:   name,
			Photo:  photo,

			Score: int(score * 100),

			TopBlock:       blockNames[best],
			SecondaryBlock: blockNames[second],

			CompatibilityType: compatibilityType(blockNames[best]),

			Explanation: explainMatch(blockNames[best], score),

			PotentialFriction: frictionExplanation(blockNames[worst]),

			CompatibilitySummary: compatibilitySummary(
				blockNames[best],
				blockNames[second],
				blockNames[worst],
			),

			RelationshipStability: stabilityLabel,

			BlockScores: blockScores,

			Role: role,
		}

		results = append(results, c)
	}

	return results, nil
}

func (r *Repository) GetTeamFit(ctx context.Context, teamID int64, userID int64) (Candidate, error) {

	query := `
	SELECT
		u.id,
		u.name,
		u.photo_url,
		p.effective_vector,
		t.effective_vector
	FROM users u
	JOIN user_mcp_profiles p ON u.id = p.user_id
	JOIN team_profiles t ON t.team_id = $1
	WHERE u.id = $2
	`

	var id int64
	var name string
	var photo string
	var vec pgvector.Vector
	var teamVec pgvector.Vector

	err := r.db.QueryRow(ctx, query, teamID, userID).
		Scan(&id, &name, &photo, &vec, &teamVec)

	if err != nil {
		return Candidate{}, err
	}

	source := vec.Slice()
	target := teamVec.Slice()

	blockSim := make([]float64, 8)
	blockScores := make(map[string]int)

	best := 0
	second := 1
	worst := 0

	for block := 0; block < 8; block++ {

		start := block * 8
		end := start + 8

		score := blockSimilarity(source, target, start, end)

		blockSim[block] = score
		blockScores[blockNames[block]] = int(score * 100)
	}

	for i := 0; i < len(blockSim); i++ {

		if blockSim[i] > blockSim[best] {
			second = best
			best = i
		}

		if blockSim[i] < blockSim[worst] {
			worst = i
		}
	}

	score := (blockSim[best] + blockSim[second]) / 2

	stabilityWeights := []float64{
		0.25,
		0.10,
		0.20,
		0.20,
		0.05,
		0.05,
		0.10,
		0.05,
	}

	stability := 0.0

	for i := 0; i < 8; i++ {
		stability += blockSim[i] * stabilityWeights[i]
	}

	stabilityLabel := StabilityScore(stability)

	role := detectRole(blockScores)

	result := Candidate{

		UserID: id,
		Name:   name,
		Photo:  photo,

		Score: int(score * 100),

		TopBlock:       blockNames[best],
		SecondaryBlock: blockNames[second],

		CompatibilityType: compatibilityType(blockNames[best]),

		Explanation: explainMatch(blockNames[best], score),

		PotentialFriction: frictionExplanation(blockNames[worst]),

		CompatibilitySummary: compatibilitySummary(
			blockNames[best],
			blockNames[second],
			blockNames[worst],
		),

		RelationshipStability: stabilityLabel,

		BlockScores: blockScores,

		Role: role,
	}

	return result, nil
}

func detectRole(blockScores map[string]int) string {

	core := blockScores["Core Values"]
	emotional := blockScores["Emotional Stability"]
	cognitive := blockScores["Cognitive Style"]
	ambition := blockScores["Ambition"]
	growth := blockScores["Growth"]
	intimacy := blockScores["Intimacy"]
	lifestyle := blockScores["Lifestyle Rhythm"]
	relational := blockScores["Relational Architecture"]

	switch {

	case ambition > 90 && core > 85:
		return "Leader"

	case cognitive > 90 && growth > 85:
		return "Strategist"

	case growth > 90 && ambition > 80:
		return "Builder"

	case lifestyle > 90 && growth > 85:
		return "Explorer"

	case relational > 90 && intimacy > 85:
		return "Connector"

	case emotional > 90 && core > 85:
		return "Stabilizer"

	case cognitive > 85 && ambition > 85:
		return "Visionary"

	case intimacy > 90 && emotional > 85:
		return "Caregiver"

	default:
		return "Balanced"
	}
}

func cosineDistance(a, b []float32) float64 {

	var dot float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {

		av := float64(a[i])
		bv := float64(b[i])

		dot += av * bv
		normA += av * av
		normB += bv * bv
	}

	return 1 - (dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}
