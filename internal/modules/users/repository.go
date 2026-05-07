package users

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/jackc/pgx/v5"
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
		COALESCE(cu.company_id, 0) AS company_id,
		u.name
	FROM users u
	LEFT JOIN company_users cu ON cu.user_id = u.id
	WHERE u.email = $1
	ORDER BY cu.is_default DESC, cu.created_at ASC
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
	SELECT id, name, COALESCE(photo_url, '')
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Company{}, err
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO companies (name, website, owner)
	VALUES ($1, $2, $3)
	RETURNING id, name, website, owner, created_at
	`

	var c models.Company

	err = tx.QueryRow(ctx, query, name, website, owner).
		Scan(&c.ID, &c.Name, &c.Website, &c.Owner, &c.CreatedAt)
	if err != nil {
		return models.Company{}, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE company_users
		SET is_default = FALSE
		WHERE user_id = $1
	`, owner)
	if err != nil {
		return models.Company{}, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO company_users (company_id, user_id, role, is_default)
		VALUES ($1, $2, 'owner', TRUE)
	`, c.ID, owner)
	if err != nil {
		return models.Company{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Company{}, err
	}

	return c, nil
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

func (r *Repository) GetTeamByID(
	ctx context.Context,
	companyID int64,
	teamID int64,
) (models.Team, error) {

	query := `
	SELECT id, company_id, name, description
	FROM teams
	WHERE id = $1 AND company_id = $2
	`

	var t models.Team

	err := r.db.QueryRow(ctx, query, teamID, companyID).
		Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description)

	return t, err
}

func (r *Repository) GetTeamMembers(
	ctx context.Context,
	companyID int64,
	teamID int64,
) ([]models.CompanyMember, error) {

	query := `
	SELECT cm.id, cm.company_id, cm.name, cm.email, cm.role, cm.created_at
	FROM team_members tm
	JOIN company_members cm ON cm.id = tm.company_member_id
	JOIN teams t ON t.id = tm.team_id
	WHERE tm.team_id = $1 AND t.company_id = $2
	ORDER BY cm.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, teamID, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.CompanyMember

	for rows.Next() {
		var m models.CompanyMember

		err := rows.Scan(
			&m.ID,
			&m.CompanyID,
			&m.Name,
			&m.Email,
			&m.Role,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		members = append(members, m)
	}

	return members, nil
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

func (r *Repository) RemoveTeamMember(
	ctx context.Context,
	companyID int64,
	teamID int64,
	memberID int64,
) error {

	query := `
	DELETE FROM team_members tm
	USING teams t, company_members cm
	WHERE tm.team_id = t.id
	  AND tm.company_member_id = cm.id
	  AND tm.team_id = $1
	  AND tm.company_member_id = $2
	  AND t.company_id = $3
	  AND cm.company_id = $3
	`

	_, err := r.db.Exec(ctx, query, teamID, memberID, companyID)

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
	candidateID *int64,
	companyMemberID *int64,
	token string,
	expiresAt time.Time,
) error {
	var entityType string
	var entityID int64

	// Asumimos que el service ya validó correctamente
	if candidateID != nil {
		entityType = "candidate"
		entityID = *candidateID
	} else {
		entityType = "member" // <-- markistente con DB
		entityID = *companyMemberID
	}

	query := `
	INSERT INTO mcp_tokens (
		entity_type,
		entity_id,
		token,
		expires_at
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entityType,
		entityID,
		token,
		expiresAt.UTC(), // <-- importante
	)

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
func (r *Repository) GetDefaultCompanyIDForUser(ctx context.Context, userID int64) (int64, error) {

	query := `
	SELECT company_id
	FROM company_users
	WHERE user_id = $1 AND is_default = TRUE
	LIMIT 1
	`

	var companyID int64

	err := r.db.QueryRow(ctx, query, userID).Scan(&companyID)
	return companyID, err
}

func (r *Repository) GetDefaultCompanyForUser(ctx context.Context, userID int64) (*models.Company, error) {

	query := `
	SELECT c.id, c.name, c.website, c.owner, c.created_at
	FROM companies c
	JOIN company_users cu ON cu.company_id = c.id
	WHERE cu.user_id = $1 AND cu.is_default = TRUE
	LIMIT 1
	`

	var company models.Company

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&company.ID,
		&company.Name,
		&company.Website,
		&company.Owner,
		&company.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}

func (r *Repository) GetCompanyForUser(
	ctx context.Context,
	userID int64,
	companyID int64,
) (*models.Company, error) {

	query := `
	SELECT c.id, c.name, c.website, c.owner, c.created_at
	FROM companies c
	JOIN company_users cu ON cu.company_id = c.id
	WHERE cu.user_id = $1 AND c.id = $2
	LIMIT 1
	`

	var company models.Company

	err := r.db.QueryRow(ctx, query, userID, companyID).Scan(
		&company.ID,
		&company.Name,
		&company.Website,
		&company.Owner,
		&company.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}

func (r *Repository) ListCompaniesForUser(
	ctx context.Context,
	userID int64,
) ([]models.CompanyMembership, error) {

	query := `
	SELECT
		c.id,
		c.name,
		c.website,
		c.owner,
		c.created_at,
		cu.role,
		cu.is_default
	FROM companies c
	JOIN company_users cu ON cu.company_id = c.id
	WHERE cu.user_id = $1
	ORDER BY cu.is_default DESC, c.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []models.CompanyMembership

	for rows.Next() {
		var membership models.CompanyMembership

		err := rows.Scan(
			&membership.ID,
			&membership.Name,
			&membership.Website,
			&membership.Owner,
			&membership.CreatedAt,
			&membership.Role,
			&membership.IsDefault,
		)
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, rows.Err()
}

func (r *Repository) SetDefaultCompanyForUser(
	ctx context.Context,
	userID int64,
	companyID int64,
) (*models.CompanyMembership, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
	SELECT
		c.id,
		c.name,
		c.website,
		c.owner,
		c.created_at,
		cu.role,
		cu.is_default
	FROM companies c
	JOIN company_users cu ON cu.company_id = c.id
	WHERE cu.user_id = $1 AND c.id = $2
	LIMIT 1
	`

	var membership models.CompanyMembership

	err = tx.QueryRow(ctx, query, userID, companyID).Scan(
		&membership.ID,
		&membership.Name,
		&membership.Website,
		&membership.Owner,
		&membership.CreatedAt,
		&membership.Role,
		&membership.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE company_users
		SET is_default = FALSE
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE company_users
		SET is_default = TRUE
		WHERE user_id = $1 AND company_id = $2
	`, userID, companyID)
	if err != nil {
		return nil, err
	}

	membership.IsDefault = true

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &membership, nil
}

func (r *Repository) ConsumeToken(
	ctx context.Context,
	token string,
) (*models.MCPToken, error) {

	query := `
	UPDATE mcp_tokens
	SET used = TRUE
	WHERE token = $1
	  AND used = FALSE
	  AND expires_at > NOW()
	RETURNING id, entity_type, entity_id, token, expires_at, used
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
func (r *Repository) GetCandidatesByCompany(ctx context.Context, companyID int64) ([]models.Candidate2, error) {

	query := `
	SELECT id, name, COALESCE(photo_url, ''), email, status, created_at
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

func (r *Repository) GetCompanyMembers(ctx context.Context, companyID int64) ([]models.CompanyMember, error) {

	query := `
	SELECT id, name, email, role, created_at
	FROM company_members
	WHERE company_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.CompanyMember

	for rows.Next() {

		var m models.CompanyMember

		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Email,
			&m.Role,
			&m.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		members = append(members, m)
	}

	return members, nil
}
func (r *Repository) GetCandidates(ctx context.Context, companyID int64) ([]models.Candidate2, error) {

	query := `
	SELECT id, name, email, COALESCE(photo_url, ''), status, created_at
	FROM candidates
	WHERE company_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Candidate2

	for rows.Next() {

		var c models.Candidate2

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.PhotoURL,
			&c.Status,
			&c.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, nil
}
func (r *Repository) GetValidToken(
	ctx context.Context,
	token string,
) (*models.MCPToken, error) {

	query := `
	SELECT id, entity_type, entity_id, token, expires_at, used
	FROM mcp_tokens
	WHERE token = $1
	  AND used = FALSE
	  AND expires_at > NOW()
	LIMIT 1
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
