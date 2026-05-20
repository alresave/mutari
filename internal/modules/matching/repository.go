package matching

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/alresave/mutari/internal/modules/mcp"
	"github.com/alresave/mutari/internal/modules/models"
	"github.com/alresave/mutari/internal/modules/users"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type Repository struct {
	db    *pgxpool.Pool
	users *users.Repository
	mcp   *mcp.Repository
}

// NewRepository creates the matching repository and its nested module repositories.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db,
		users: users.NewRepository(db),
		mcp:   mcp.NewRepository(db),
	}
}

// Users exposes user/company persistence through its interface boundary.
func (r *Repository) Users() users.RepositoryStore {
	return r.users
}

// MCP exposes MCP persistence through its interface boundary.
func (r *Repository) MCP() mcp.RepositoryStore {
	return r.mcp
}

func blockSimilarity(a, b []float32, start, end int) float64 {
	if start < 0 || end > len(a) || end > len(b) || start >= end {
		return 0
	}

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

// StabilityScore turns a normalized compatibility score into a user-facing label.
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

	case "Professional Values", "Core Values":
		return "Values Match"

	case "Collaboration Depth", "Relational Architecture":
		return "Relationship Structure Match"

	case "Operational Rhythm", "Lifestyle Rhythm":
		return "Lifestyle Match"

	case "Team Positioning", "Growth":
		return "Growth Alignment"

	case "Cognitive Processing", "Cognitive Style":
		return "Thinking Style Match"

	case "Learning Velocity", "Ambition":
		return "Ambition Alignment"

	case "Emotional Resilience", "Emotional Stability":
		return "Emotional Compatibility"

	case "Execution Drive", "Trust & Collaboration":
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

func (r *Repository) GetMemberMatches(ctx context.Context, memberID int64) ([]models.Candidate, error) {

	query := `
	SELECT
	cm.id,
	cm.name,
	'' AS photo_url,
	p.effective_vector
	FROM company_member_mcp_profiles p
	JOIN company_members cm ON cm.id = p.company_member_id
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	source, ok := vectors[memberID]
	if !ok {
		return nil, errors.New("member vector not found")
	}

	var results []models.Candidate

	for id, vec := range vectors {

		if id == memberID {
			continue
		}

		blockScores := make(map[string]int)
		blockSim := make([]float64, models.MCPBlockCount)

		best := 0
		second := 1
		worst := 0

		for block := 0; block < models.MCPBlockCount; block++ {

			start := block * models.MCPDimensionsPerBlock
			end := start + models.MCPDimensionsPerBlock

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

		// TODO: make role weights dynamic when member role preferences are stored.
		weights := models.RoleWeights["engineer"]

		score := 0.0

		for i := 0; i < models.MCPBlockCount; i++ {
			score += blockSim[i] * weights[i]
		}

		c := models.Candidate{

			MemberID: id,
			Name:     names[id],
			Photo:    photos[id],

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
		if errors.Is(err, pgx.ErrNoRows) {
			// Lazily build the team profile when a hiring profile is requested
			// before the aggregate vector exists.
			if buildErr := r.BuildTeamVector(ctx, teamID); buildErr != nil {
				return nil, buildErr
			}

			err = r.db.QueryRow(ctx, query, teamID).Scan(&vec)
		}
	}

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
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) BuildTeamVector(ctx context.Context, teamID int64) error {

	query := `
 SELECT p.effective_vector
 FROM team_members tm
 JOIN company_member_mcp_profiles p ON p.company_member_id = tm.company_member_id
 WHERE tm.team_id = $1
 `

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Team profiles are the centroid of all available member MCP vectors.
	vector := make([]float32, models.MCPVectorSize)
	count := 0

	for rows.Next() {

		var v pgvector.Vector

		if err := rows.Scan(&v); err != nil {
			return err
		}

		values := v.Slice()

		if len(values) != models.MCPVectorSize {
			return fmt.Errorf("member vector has %d dimensions, expected %d", len(values), models.MCPVectorSize)
		}

		for i := 0; i < models.MCPVectorSize; i++ {
			vector[i] += values[i]
		}

		count++
	}

	if count == 0 {
		return nil
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for i := 0; i < models.MCPVectorSize; i++ {
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) GetTeamMatches(ctx context.Context, teamID int64) ([]models.Candidate, error) {

	query := `
	SELECT
	c.id,
	c.name,
	COALESCE(c.photo_url, ''),
	p.effective_vector,
	t.effective_vector
	FROM team_profiles t
	JOIN candidate_mcp_profiles p ON TRUE
	JOIN candidates c ON c.id = p.candidate_id
	WHERE t.team_id=$1
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

		blockSim := make([]float64, models.MCPBlockCount)
		blockScores := make(map[string]int)

		best := 0
		second := 1
		worst := 0

		for block := 0; block < models.MCPBlockCount; block++ {

			start := block * models.MCPDimensionsPerBlock
			end := start + models.MCPDimensionsPerBlock

			score := blockSimilarity(source, target, start, end)

			blockSim[block] = score
			blockScores[models.BlockNames[block]] = int(score * 100)
		}

		for i := 0; i < models.MCPBlockCount; i++ {

			if blockSim[i] > blockSim[best] {
				second = best
				best = i
			}

			if blockSim[i] < blockSim[worst] {
				worst = i
			}
		}

		score := (blockSim[best] + blockSim[second]) / 2

		// Relationship stability gives heavier influence to values, emotional
		// fit, and relational architecture than to skill-adjacent traits.
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

		for i := 0; i < models.MCPBlockCount; i++ {
			stability += blockSim[i] * stabilityWeights[i]
		}

		stabilityLabel := StabilityScore(stability)

		role := predictRole(blockSim)

		c := models.Candidate{

			MemberID: id,
			Name:     name,
			Photo:    photo,

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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) GetTeamMemberFit(ctx context.Context, teamID int64, memberID int64) (models.Candidate, error) {

	query := `
	SELECT
		cm.id,
		cm.name,
		'' AS photo_url,
		p.effective_vector,
		t.effective_vector
	FROM company_members cm
	JOIN company_member_mcp_profiles p ON cm.id = p.company_member_id
	JOIN team_profiles t ON t.team_id = $1
	WHERE cm.id = $2
	`

	var id int64
	var name string
	var photo string
	var vec pgvector.Vector
	var teamVec pgvector.Vector

	err := r.db.QueryRow(ctx, query, teamID, memberID).
		Scan(&id, &name, &photo, &vec, &teamVec)

	if err != nil {
		return models.Candidate{}, err
	}

	source := vec.Slice()
	target := teamVec.Slice()

	blockSim := make([]float64, models.MCPBlockCount)
	blockScores := make(map[string]int)

	best := 0
	second := 1
	worst := 0

	for block := 0; block < models.MCPBlockCount; block++ {

		start := block * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock

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

	for i := 0; i < models.MCPBlockCount; i++ {
		stability += blockSim[i] * stabilityWeights[i]
	}

	stabilityLabel := StabilityScore(stability)

	role := predictRole(blockSim)

	result := models.Candidate{

		MemberID: id,
		Name:     name,
		Photo:    photo,

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
	JOIN company_member_mcp_profiles p ON p.company_member_id = tm.company_member_id
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

		values := vec.Slice()
		if len(values) != models.MCPVectorSize {
			return nil, "", fmt.Errorf("member vector has %d dimensions, expected %d", len(values), models.MCPVectorSize)
		}

		vectors = append(vectors, values)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	if len(vectors) == 0 {
		return nil, "", errors.New("team has no members")
	}

	// Collapse each member's 64-dimensional profile into the eight MCP blocks,
	// then average those block scores across the team.
	blockScores := make(map[string]int)

	for block := 0; block < models.MCPBlockCount; block++ {

		start := block * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock

		sum := 0.0

		for _, vec := range vectors {

			local := 0.0

			for i := start; i < end; i++ {
				local += float64(vec[i])
			}

			local = local / models.MCPDimensionsPerBlock
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
 JOIN company_member_mcp_profiles p ON p.company_member_id = tm.company_member_id
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

		values := vec.Slice()
		if len(values) != models.MCPVectorSize {
			return nil, "", fmt.Errorf("member vector has %d dimensions, expected %d", len(values), models.MCPVectorSize)
		}

		vectors = append(vectors, values)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	if len(vectors) == 0 {
		return nil, "", errors.New("team has no members")
	}

	// Reuse the team block averages to identify the two weakest cultural areas.
	blockScores := make(map[string]int)

	for block := 0; block < models.MCPBlockCount; block++ {

		start := block * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock

		sum := 0.0

		for _, vec := range vectors {

			local := 0.0

			for i := start; i < end; i++ {
				local += float64(vec[i])
			}

			local = local / models.MCPDimensionsPerBlock
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

func (r *Repository) GetTeams(ctx context.Context, companyID int64) ([]models.Team, error) {

	query := `
SELECT id, name,description, company_id
FROM teams
WHERE company_id = $1
ORDER BY id
`

	rows, err := r.db.Query(ctx, query, companyID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
func (r *Repository) AddMemberToTeam(ctx context.Context, teamID int64, memberID int64) error {

	query := `
    INSERT INTO team_members (team_id, company_member_id)
    VALUES ($1,$2)
    ON CONFLICT DO NOTHING
    `

	_, err := r.db.Exec(ctx, query, teamID, memberID)

	return err
}
func (r *Repository) GetCandidates(ctx context.Context, teamID int64) ([]models.User, error) {

	query := `
    SELECT id,name,COALESCE(photo_url, '')
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
	if err := rows.Err(); err != nil {
		return nil, err
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

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

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

	// Role suggestions are based on block-level averages rather than individual
	// dimensions so small per-question noise does not dominate the label.
	blockAverage := func(block int) float32 {
		start := block * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock
		if len(vec) < end {
			return 0
		}

		var sum float32
		for i := start; i < end; i++ {
			sum += vec[i]
		}

		return sum / models.MCPDimensionsPerBlock
	}

	execution := blockAverage(0)
	cognitive := blockAverage(1)
	emotional := blockAverage(3)

	if cognitive > 0.75 {
		return "Strategist"
	}

	if execution > 0.75 {
		return "Leader"
	}

	if emotional > 0.75 {
		return "Connector"
	}

	return "Specialist"
}
func computeBlockScores(team, candidate []float32) map[string]int {

	scores := make(map[string]int)

	for block := 0; block < models.MCPBlockCount; block++ {

		start := block * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock

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

	// Load the candidate profile that will be compared against each team member.
	var candidateVec pgvector.Vector

	err := r.db.QueryRow(ctx, `
	SELECT effective_vector
	FROM candidate_mcp_profiles
	WHERE candidate_id=$1
	`, candidateID).Scan(&candidateVec)

	if err != nil {
		return models.TeamFitSimulation{}, err
	}

	candidate := candidateVec.Slice()
	if len(candidate) != models.MCPVectorSize {
		return models.TeamFitSimulation{}, fmt.Errorf("candidate vector has %d dimensions, expected %d", len(candidate), models.MCPVectorSize)
	}

	// Use individual member vectors instead of only the team centroid so the
	// simulation can account for balance and variance across the team.
	rows, err := r.db.Query(ctx, `
	SELECT p.effective_vector
	FROM team_members tm
	JOIN company_member_mcp_profiles p ON p.company_member_id = tm.company_member_id
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

		values := vec.Slice()
		if len(values) != len(candidate) {
			return models.TeamFitSimulation{}, fmt.Errorf("member vector has %d dimensions, expected %d", len(values), len(candidate))
		}

		teamVectors = append(teamVectors, values)
	}
	if err := rows.Err(); err != nil {
		return models.TeamFitSimulation{}, err
	}

	if len(teamVectors) == 0 {
		return models.TeamFitSimulation{}, errors.New("team has no members")
	}

	// Average the candidate's pairwise similarity to every current team member.
	totalSim := 0.0
	sims := make([]float64, len(teamVectors))

	for i, member := range teamVectors {

		sim := cosineSimilarity(member, candidate)

		sims[i] = sim
		totalSim += sim
	}

	teamFit := totalSim / float64(len(teamVectors))

	// Penalize uneven fit: a candidate who matches only part of the team should
	// score lower than one with consistent compatibility across members.
	var variance float64

	for _, s := range sims {
		diff := s - teamFit
		variance += diff * diff
	}

	variance /= float64(len(sims))

	balancePenalty := math.Sqrt(variance)

	// Add a small role-complementarity boost after the base fit is calculated.
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

	// Convert the blended cosine score into a 0-100 percentage-like score.
	scoreRaw :=
		(teamFit * 0.75) +
			((1 - balancePenalty) * 0.20) +
			roleBoost
	if scoreRaw < -1 {
		scoreRaw = -1
	}
	if scoreRaw > 1 {
		scoreRaw = 1
	}

	dist := math.Acos(scoreRaw) / math.Pi

	score := int((1 - dist) * 100)

	// Compute the team centroid only for the returned per-block score breakdown.
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

	if scoreForIntBlock(scores, "Learning Velocity", "Ambition") > 85 &&
		scoreForIntBlock(scores, "Cognitive Processing", "Cognitive Style") > 85 {
		return "Visionary Builders"
	}

	if scoreForIntBlock(scores, "Collaboration Depth", "Relational Architecture") > 85 {
		return "Collaborative Operators"
	}

	if scoreForIntBlock(scores, "Emotional Resilience", "Emotional Stability") > 85 {
		return "Stable Executors"
	}

	return "Balanced Team"
}
func recommendProfile(blocks []string) string {

	for _, b := range blocks {

		switch b {

		case "Emotional Resilience", "Emotional Stability":
			return "Stabilizer"

		case "Collaboration Depth", "Relational Architecture":
			return "Collaborator"

		case "Cognitive Processing", "Cognitive Style":
			return "Strategist"

		case "Learning Velocity", "Ambition":
			return "Driver"

		case "Team Positioning", "Growth":
			return "Learner"

		case "Execution Drive", "Trust & Collaboration":
			return "Connector"

		case "Operational Rhythm", "Lifestyle Rhythm":
			return "Balancer"

		case "Professional Values", "Core Values":
			return "Anchor"
		}
	}

	return "Generalist"
}

func scoreForIntBlock(scores map[string]int, names ...string) int {
	for _, name := range names {
		if score, ok := scores[name]; ok {
			return score
		}
	}

	return 0
}

func PredictTeamArchetype(vector []float32) string {

	if len(vector) < models.MCPVectorSize {
		return "Balanced Team"
	}

	blockAverage := func(block int) float32 {
		start := block * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock

		var sum float32
		for i := start; i < end; i++ {
			sum += vector[i]
		}

		return sum / models.MCPDimensionsPerBlock
	}

	cognitive := blockAverage(1)
	ambition := blockAverage(4)
	emotional := blockAverage(3)

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
