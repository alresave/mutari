package users

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (int64, string, int64, string, error) {

	query := `
	SELECT 
		u.id,
		u.password_hash,
		COALESCE(c.id, 0) AS company_id,
		u.name
	FROM users u
	LEFT JOIN companies c ON c.owner = u.id
	WHERE u.email = $1
	LIMIT 1
	`

	var id int64
	var hash string
	var companyID int64
	var name string

	err := r.db.QueryRow(ctx, query, email).Scan(
		&id,
		&hash,
		&companyID,
		&name,
	)

	if err != nil {
		return 0, "", 0, "", err
	}

	return id, hash, companyID, name, nil
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

func (r *Repository) CreateCompany(ctx context.Context, name, website string, owner int64) (models.Company, error) {

	query := `
	INSERT INTO companies (name, website, owner)
	VALUES ($1, $2, $3)
	RETURNING id, name, website, owner, created_at
	`

	var c models.Company

	err := r.db.QueryRow(ctx, query, name, website, owner).
		Scan(&c.ID, &c.Name, &c.Website, &c.Owner, &c.CreatedAt)

	return c, err
}

func (r *Repository) CreateCompanyMember(
	ctx context.Context,
	companyID int64,
	name, email, role string,
) (models.CompanyMember, error) {

	query := `
	INSERT INTO company_members (company_id, name, email, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id, company_id, name, email, role, created_at
	`
	var m models.CompanyMember

	err := r.db.QueryRow(ctx, query, companyID, name, email, role).
		Scan(&m.ID, &m.CompanyID, &m.Name, &m.Email, &m.Role, &m.CreatedAt)

	return m, err
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

func (r *Repository) CreateTeam(
	ctx context.Context,
	companyID int64,
	name, description string,
) (models.Team, error) {

	query := `
	INSERT INTO teams (company_id, name, description)
	VALUES ($1, $2, $3)
	RETURNING id, company_id, name, description
	`

	var t models.Team

	err := r.db.QueryRow(ctx, query, companyID, name, description).
		Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description)

	return t, err
}
func (r *Repository) AddTeamMember(
	ctx context.Context,
	teamID int64,
	memberID int64,
) error {

	query := `
	INSERT INTO team_members (team_id, company_member_id)
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, teamID, memberID)

	return err
}

func (r *Repository) TeamExists(ctx context.Context, teamID int64) (bool, error) {

	query := `
	SELECT EXISTS (
		SELECT 1 FROM teams WHERE id = $1
	)
	`

	var exists bool
	log.Printf("Checking if team %d exists", teamID)

	err := r.db.QueryRow(ctx, query, teamID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) TeamExistsForCompany(
	ctx context.Context,
	teamID int64,
	companyID int64,
) (bool, error) {

	query := `
	SELECT EXISTS (
		SELECT 1 FROM teams
		WHERE id = $1 AND company_id = $2
	)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, teamID, companyID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) CompanyMemberExistsForCompany(
	ctx context.Context,
	memberID int64,
	companyID int64,
) (bool, error) {

	query := `
	SELECT EXISTS (
		SELECT 1
		FROM company_members
		WHERE id = $1 AND company_id = $2
	)
	`

	var exists bool

	err := r.db.QueryRow(ctx, query, memberID, companyID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (r *Repository) CreateToken(
	ctx context.Context,
	candidateID int64,
	token string,
	expiresAt time.Time,
) error {

	query := `
	INSERT INTO mcp_tokens (candidate_id, token, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, candidateID, token, expiresAt)
	return err
}
func (r *Repository) CreateCandidate(
	ctx context.Context,
	companyID int64,
	name string,
	email string,
) (int64, error) {

	query := `
	INSERT INTO candidates (company_id, name, email)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var id int64

	err := r.db.QueryRow(ctx, query, companyID, name, email).Scan(&id)
	return id, err
}

func (r *Repository) GetToken(ctx context.Context, token string) (*models.MCPToken, error) {

	query := `
	SELECT id, entity_type, entity_id, token, expires_at, used
	FROM mcp_tokens
	WHERE token = $1
	`

	var t models.MCPToken

	err := r.db.QueryRow(ctx, query, token).Scan(
		&t.ID,
		&t.EntityType,
		&t.EntityID,
		&t.Token,
		&t.ExpiresAt,
		&t.Used,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
func (r *Repository) GetCompanyByOwner(ctx context.Context, userID int64) (int64, error) {

	query := `
	SELECT id FROM companies
	WHERE owner = $1
	LIMIT 1
	`

	var companyID int64

	err := r.db.QueryRow(ctx, query, userID).Scan(&companyID)
	return companyID, err
}
func (r *Repository) InsertCandidateAnswers(
	ctx context.Context,
	candidateID int64,
	answers []models.AnswerInput,
) error {

	query := `
	INSERT INTO candidate_mcp_answers (candidate_id, question_id, value)
	VALUES ($1, $2, $3)
	`

	for _, a := range answers {
		_, err := r.db.Exec(ctx, query,
			candidateID,
			a.QuestionID,
			a.Value,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *Repository) InsertCompanyMemberAnswers(
	ctx context.Context,
	memberID int64,
	answers []models.AnswerInput,
) error {

	query := `
	INSERT INTO company_member_mcp_answers (company_member_id, question_id, value)
	VALUES ($1, $2, $3)
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

func (r *Repository) UpsertCandidateProfile(
	ctx context.Context,
	candidateID int64,
	vector []float32,
) error {

	query := `
	INSERT INTO candidate_mcp_profiles (candidate_id, effective_vector)
	VALUES ($1, $2)
	ON CONFLICT (candidate_id)
	DO UPDATE SET effective_vector = EXCLUDED.effective_vector
	`

	_, err := r.db.Exec(ctx, query, candidateID, vector)
	return err
}
func (r *Repository) UpsertCompanyMemberProfile(
	ctx context.Context,
	memberID int64,
	vector []float32,
) error {

	query := `
	INSERT INTO company_member_profiles (company_member_id, effective_vector)
	VALUES ($1, $2)
	ON CONFLICT (company_member_id)
	DO UPDATE SET effective_vector = EXCLUDED.effective_vector
	`

	_, err := r.db.Exec(ctx, query, memberID, vector)
	return err
}

func (r *Repository) MarkTokenUsed(ctx context.Context, token string) error {

	query := `
	UPDATE mcp_tokens
	SET used = true
	WHERE token = $1
	`

	_, err := r.db.Exec(ctx, query, token)
	return err
}
func BuildBlockVector(
	questions []models.MCPQuestion,
	answers []models.AnswerInput,
) []float32 {

	answerMap := make(map[int64]float32)
	for _, a := range answers {
		answerMap[a.QuestionID] = float32(a.Value)
	}

	blockSum := make(map[int64]float32)
	blockCount := make(map[int64]int)

	// 🔹 registrar TODOS los bloques
	blockSet := make(map[int64]struct{})

	for _, q := range questions {

		blockSet[q.Block] = struct{}{}

		if val, ok := answerMap[q.ID]; ok {
			blockSum[q.Block] += val
			blockCount[q.Block]++
		}
	}

	// 🔹 ordenar bloques
	var blocks []int64
	for b := range blockSet {
		blocks = append(blocks, b)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i] < blocks[j]
	})

	vector := make([]float32, len(blocks))

	for i, b := range blocks {

		count := blockCount[b]

		if count == 0 {
			vector[i] = 0
			continue
		}

		avg := blockSum[b] / float32(count)
		vector[i] = avg / 5.0
	}

	return vector
}
func (r *Repository) GetCandidatesByCompany(ctx context.Context, companyID int64) ([]models.Candidate2, error) {

	query := `
	SELECT id, name, photo_url, email, status, created_at
	FROM candidates
	WHERE company_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []models.Candidate2

	for rows.Next() {
		var c models.Candidate2

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.PhotoURL,
			&c.Email,
			&c.Status,
			&c.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		candidates = append(candidates, c)
	}

	return candidates, nil
}
