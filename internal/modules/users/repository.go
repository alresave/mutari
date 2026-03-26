package users

import (
	"context"
	"log"
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
    SELECT id, password_hash, company_id, name
    FROM users
    WHERE email = $1
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
