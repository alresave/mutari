package matching

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/alresave/mutari/internal/modules/mcp"
	"github.com/alresave/mutari/internal/modules/models"
	"github.com/alresave/mutari/internal/modules/users"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type Repository struct {
	db    *pgxpool.Pool
	users *users.Repository
	mcp   *mcp.Repository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db,
		users: users.NewRepository(db),
		mcp:   mcp.NewRepository(db),
	}
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

	case "Trust & Collaboration ":
		return "Bonding Compatibility"

	default:
		return "General Compatibility"
	}
}

func frictionExplanation(block string) string {

	return models.BlockDifferences[block]
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

func (r *Repository) GetMCPMatches(ctx context.Context, userID int64) ([]models.Candidate, error) {

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

	var results []models.Candidate

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
			blockScores[models.BlockNames[block]] = int(score * 100)
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

		weights := models.RoleWeights["engineer"] // luego dinámico

		score := 0.0

		for i := 0; i < 8; i++ {
			score += blockSim[i] * weights[i]
		}

		c := models.Candidate{

			UserID: id,
			Name:   names[id],
			Photo:  photos[id],

			Score: int(score * 100),

			TopBlock:       models.BlockNames[best],
			SecondaryBlock: models.BlockNames[second],

			CompatibilityType: compatibilityType(models.BlockNames[best]),

			Explanation: explainMatch(models.BlockNames[best], score),

			PotentialFriction: frictionExplanation(models.BlockNames[worst]),

			CompatibilitySummary: compatibilitySummary(
				models.BlockNames[best],
				models.BlockNames[second],
				models.BlockNames[worst],
			),

			RelationshipStability: StabilityScore(score),
			Role:                  predictRole(blockSim),

			BlockScores: blockScores,
		}

		results = append(results, c)
	}

	return results, nil
}

func (r *Repository) GetTeamVector(ctx context.Context, teamID int64) (*models.TeamVector, error) {

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

	vector := vec.Slice()

	var score float64

	for _, v := range vector {
		score += float64(v)
	}

	score = score / float64(len(vector)) * 100

	archetype := PredictTeamArchetype(vector)

	return &models.TeamVector{
		Vector:    vector,
		TeamScore: score,
		Archetype: archetype,
	}, nil
}

func (r *Repository) GetTeamHiringProfile(ctx context.Context, teamID int64) ([]float32, error) {

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

	return vec.Slice(), nil
}

func (r *Repository) RebuildAllTeams(ctx context.Context) error {

	query := `SELECT id FROM teams`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		var teamID int64

		if err := rows.Scan(&teamID); err != nil {
			return err
		}

		err := r.BuildTeamVector(ctx, teamID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) BuildTeamVector(ctx context.Context, teamID int64) error {

	query := `
 SELECT p.effective_vector
 FROM team_members tm
 JOIN user_mcp_profiles p ON p.user_id = tm.user_id
 WHERE tm.team_id = $1
 `

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return err
	}
	defer rows.Close()

	vector := make([]float32, 64)
	count := 0

	for rows.Next() {

		var v pgvector.Vector

		if err := rows.Scan(&v); err != nil {
			return err
		}

		values := v.Slice()

		for i := 0; i < 64; i++ {
			vector[i] += values[i]
		}

		count++
	}

	if count == 0 {
		return nil
	}

	for i := 0; i < 64; i++ {
		vector[i] = vector[i] / float32(count)
	}

	pgVector := pgvector.NewVector(vector)

	queryInsert := `
 INSERT INTO team_profiles (team_id,effective_vector)
 VALUES ($1,$2)
 ON CONFLICT(team_id)
 DO UPDATE SET
 effective_vector = EXCLUDED.effective_vector,
 updated_at = NOW()
 `

	_, err = r.db.Exec(ctx, queryInsert, teamID, pgVector)

	return err
}
func (r *Repository) GetAllTeamIDs(ctx context.Context) ([]int64, error) {

	query := `SELECT id FROM teams`

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

func (r *Repository) GetTeamMatches(ctx context.Context, teamID int64) ([]models.Candidate, error) {

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

	var results []models.Candidate

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
			blockScores[models.BlockNames[block]] = int(score * 100)
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

		role := predictRole(stabilityWeights)

		c := models.Candidate{

			UserID: id,
			Name:   name,
			Photo:  photo,

			Score: int(score * 100),

			TopBlock:       models.BlockNames[best],
			SecondaryBlock: models.BlockNames[second],

			CompatibilityType: compatibilityType(models.BlockNames[best]),

			Explanation: explainMatch(models.BlockNames[best], score),

			PotentialFriction: frictionExplanation(models.BlockNames[worst]),

			CompatibilitySummary: compatibilitySummary(
				models.BlockNames[best],
				models.BlockNames[second],
				models.BlockNames[worst],
			),

			RelationshipStability: stabilityLabel,

			BlockScores: blockScores,

			Role: role,
		}

		results = append(results, c)
	}

	return results, nil
}

func (r *Repository) GetTeamFit(ctx context.Context, teamID int64, userID int64) (models.Candidate, error) {

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
		return models.Candidate{}, err
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
		blockScores[models.BlockNames[block]] = int(score * 100)
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

	role := predictRole(stabilityWeights)

	result := models.Candidate{

		UserID: id,
		Name:   name,
		Photo:  photo,

		Score: int(score * 100),

		TopBlock:       models.BlockNames[best],
		SecondaryBlock: models.BlockNames[second],

		CompatibilityType: compatibilityType(models.BlockNames[best]),

		Explanation: explainMatch(models.BlockNames[best], score),

		PotentialFriction: frictionExplanation(models.BlockNames[worst]),

		CompatibilitySummary: compatibilitySummary(
			models.BlockNames[best],
			models.BlockNames[second],
			models.BlockNames[worst],
		),

		RelationshipStability: stabilityLabel,

		BlockScores: blockScores,

		Role: role,
	}

	return result, nil
}

func (r *Repository) GetTeamCulture(ctx context.Context, teamID int64) (map[string]int, string, error) {

	query := `
	SELECT p.effective_vector
	FROM team_members tm
	JOIN user_mcp_profiles p ON p.user_id = tm.user_id
	WHERE tm.team_id = $1
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var vectors [][]float32

	for rows.Next() {

		var vec pgvector.Vector

		if err := rows.Scan(&vec); err != nil {
			return nil, "", err
		}

		vectors = append(vectors, vec.Slice())
	}

	if len(vectors) == 0 {
		return nil, "", errors.New("team has no members")
	}

	blockScores := make(map[string]int)

	for block := 0; block < 8; block++ {

		start := block * 8
		end := start + 8

		sum := 0.0

		for _, vec := range vectors {

			local := 0.0

			for i := start; i < end; i++ {
				local += float64(vec[i])
			}

			local = local / 8
			sum += local
		}

		avg := sum / float64(len(vectors))

		blockScores[models.BlockNames[block]] = int(avg * 100)
	}

	archetype := detectTeamArchetype(blockScores)

	return blockScores, archetype, nil
}
func (r *Repository) GetTeamGaps(ctx context.Context, teamID int64) ([]string, string, error) {

	query := `
 SELECT p.effective_vector
 FROM team_members tm
 JOIN user_mcp_profiles p ON p.user_id = tm.user_id
 WHERE tm.team_id = $1
 `

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var vectors [][]float32

	for rows.Next() {

		var vec pgvector.Vector

		if err := rows.Scan(&vec); err != nil {
			return nil, "", err
		}

		vectors = append(vectors, vec.Slice())
	}

	if len(vectors) == 0 {
		return nil, "", errors.New("team has no members")
	}

	blockScores := make(map[string]int)

	for block := 0; block < 8; block++ {

		start := block * 8
		end := start + 8

		sum := 0.0

		for _, vec := range vectors {

			local := 0.0

			for i := start; i < end; i++ {
				local += float64(vec[i])
			}

			local = local / 8
			sum += local
		}

		avg := sum / float64(len(vectors))

		blockScores[models.BlockNames[block]] = int(avg * 100)
	}

	type pair struct {
		block string
		score int
	}

	var list []pair

	for k, v := range blockScores {
		list = append(list, pair{k, v})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].score < list[j].score
	})

	weak := []string{
		list[0].block,
		list[1].block,
	}

	profile := recommendProfile(weak)

	return weak, profile, nil
}

func (r *Repository) GetTeams(ctx context.Context) ([]models.Team, error) {

	query := `
SELECT id, name,description, company_id
FROM teams
ORDER BY id
`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team

	for rows.Next() {

		var t models.Team

		err := rows.Scan(
			&t.ID,
			&t.Name,

			&t.Description,
			&t.CompanyID,
		)

		if err != nil {
			return nil, err
		}

		teams = append(teams, t)

	}

	return teams, nil
}
func (r *Repository) AddUserToTeam(ctx context.Context, teamID int64, userID int64) error {

	query := `
    INSERT INTO team_members (team_id, user_id)
    VALUES ($1,$2)
    ON CONFLICT DO NOTHING
    `

	_, err := r.db.Exec(ctx, query, teamID, userID)

	return err
}
func (r *Repository) GetCandidates(ctx context.Context, teamID int64) ([]models.User, error) {

	query := `
    SELECT id,name,photo_url
    FROM users
    WHERE id NOT IN (
        SELECT user_id
        FROM team_members
        WHERE team_id=$1
    )
    `

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var u models.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Photo); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
func (r *Repository) GetTeamProfile(ctx context.Context, teamID int64) ([]float32, error) {

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

	return vec.Slice(), nil
}

func detectRole(blockScores map[string]int) string {

	core := blockScores["Core Values"]
	emotional := blockScores["Emotional Stability"]
	cognitive := blockScores["Cognitive Style"]
	ambition := blockScores["Ambition"]
	growth := blockScores["Growth"]
	TrustCollaboration := blockScores["Trust & Collaboration "]
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

	case relational > 90 && TrustCollaboration > 85:
		return "Connector"

	case emotional > 90 && core > 85:
		return "Stabilizer"

	case cognitive > 85 && ambition > 85:
		return "Visionary"

	case TrustCollaboration > 90 && emotional > 85:
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
func cosineSimilarity(a, b []float32) float64 {

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

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
func suggestRole(vec []float32) string {

	cognitive := vec[56]
	ambition := vec[32]
	emotional := vec[16]

	if cognitive > 0.75 {
		return "Strategist"
	}

	if ambition > 0.75 {
		return "Leader"
	}

	if emotional > 0.75 {
		return "Connector"
	}

	return "Specialist"
}
func computeBlockScores(team, candidate []float32) map[string]int {

	scores := make(map[string]int)

	for block := 0; block < 8; block++ {

		start := block * 8
		end := start + 8

		score := blockSimilarity(team, candidate, start, end)

		scores[models.BlockNames[block]] = int(score * 100)
	}

	return scores
}

func (r *Repository) SimulateTeamFit(
	ctx context.Context,
	teamID int64,
	candidateID int64,
) (models.TeamFitSimulation, error) {

	// -----------------------------
	// Candidate vector
	// -----------------------------

	var candidateVec pgvector.Vector

	err := r.db.QueryRow(ctx, `
	SELECT effective_vector
	FROM user_mcp_profiles
	WHERE user_id=$1
	`, candidateID).Scan(&candidateVec)

	if err != nil {
		return models.TeamFitSimulation{}, err
	}

	candidate := candidateVec.Slice()

	// -----------------------------
	// Team member vectors
	// -----------------------------

	rows, err := r.db.Query(ctx, `
	SELECT p.effective_vector
	FROM team_members tm
	JOIN user_mcp_profiles p ON p.user_id = tm.user_id
	WHERE tm.team_id=$1
	`, teamID)

	if err != nil {
		return models.TeamFitSimulation{}, err
	}

	defer rows.Close()

	var teamVectors [][]float32

	for rows.Next() {

		var vec pgvector.Vector

		if err := rows.Scan(&vec); err != nil {
			return models.TeamFitSimulation{}, err
		}

		teamVectors = append(teamVectors, vec.Slice())
	}

	if len(teamVectors) == 0 {
		return models.TeamFitSimulation{}, errors.New("team has no members")
	}

	// -----------------------------
	// Candidate vs each team member
	// -----------------------------

	totalSim := 0.0
	sims := make([]float64, len(teamVectors))

	for i, member := range teamVectors {

		sim := cosineSimilarity(member, candidate)

		sims[i] = sim
		totalSim += sim
	}

	teamFit := totalSim / float64(len(teamVectors))

	// -----------------------------
	// Team balance penalty
	// -----------------------------

	var variance float64

	for _, s := range sims {
		diff := s - teamFit
		variance += diff * diff
	}

	variance /= float64(len(sims))

	balancePenalty := math.Sqrt(variance)

	// -----------------------------
	// Role complementarity
	// -----------------------------

	role := suggestRole(candidate)

	roleBoost := 0.0

	switch role {

	case "Strategist":
		roleBoost = 0.05

	case "Leader":
		roleBoost = 0.05

	case "Connector":
		roleBoost = 0.04

	default:
		roleBoost = 0.02
	}

	// -----------------------------
	// Final score
	// -----------------------------

	scoreRaw :=
		(teamFit * 0.75) +
			((1 - balancePenalty) * 0.20) +
			roleBoost

	dist := math.Acos(scoreRaw) / math.Pi

	score := int((1 - dist) * 100)

	// -----------------------------
	// Compute team average vector
	// (for block scores)
	// -----------------------------

	dim := len(candidate)

	avg := make([]float32, dim)

	for _, v := range teamVectors {
		for i := 0; i < dim; i++ {
			avg[i] += v[i]
		}
	}

	for i := 0; i < dim; i++ {
		avg[i] /= float32(len(teamVectors))
	}

	blockScores := computeBlockScores(avg, candidate)

	return models.TeamFitSimulation{
		TeamID:      teamID,
		CandidateID: candidateID,
		Score:       score,
		Role:        role,
		BlockScores: blockScores,
	}, nil
}

func detectTeamArchetype(scores map[string]int) string {

	if scores["Ambition"] > 85 && scores["Cognitive Style"] > 85 {
		return "Visionary Builders"
	}

	if scores["Relational Architecture"] > 85 {
		return "Collaborative Operators"
	}

	if scores["Emotional Stability"] > 85 {
		return "Stable Executors"
	}

	return "Balanced Team"
}
func recommendProfile(blocks []string) string {

	for _, b := range blocks {

		switch b {

		case "Emotional Stability":
			return "Stabilizer"

		case "Relational Architecture":
			return "Collaborator"

		case "Cognitive Style":
			return "Strategist"

		case "Ambition":
			return "Driver"

		case "Growth":
			return "Learner"

		case "Trust & Collaboration ":
			return "Connector"

		case "Lifestyle Rhythm":
			return "Balancer"

		case "Core Values":
			return "Anchor"
		}
	}

	return "Generalist"
}
func PredictTeamArchetype(vector []float32) string {

	if len(vector) < 8 {
		return "Balanced Team"
	}

	cognitive := vector[1]
	ambition := vector[5]
	emotional := vector[2]

	if ambition > 0.75 && cognitive > 0.75 {
		return "Visionary Team"
	}

	if ambition > 0.75 {
		return "Execution Team"
	}

	if emotional > 0.75 {
		return "Stability Team"
	}

	return "Balanced Team"
}
