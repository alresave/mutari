package users

import (
	"context"
	"time"

	"github.com/alresave/mutari/internal/modules/models"
)

// RepositoryStore defines the user/company persistence behavior required by services.
type RepositoryStore interface {
	GetUserByEmail(ctx context.Context, email string) (int64, string, int64, string, error)
	GetUserByID(ctx context.Context, userID int64) (int64, string, string, error)
	CreateUser(ctx context.Context, name, email, passwordHash string) (int64, error)
	CreateCompany(ctx context.Context, name, website string, owner int64) (models.Company, error)
	ListCompaniesForUser(ctx context.Context, userID int64) ([]models.CompanyMembership, error)
	SetDefaultCompanyForUser(ctx context.Context, userID int64, companyID int64) (*models.CompanyMembership, error)
	GetDefaultCompanyForUser(ctx context.Context, userID int64) (*models.Company, error)
	GetCompanyForUser(ctx context.Context, userID int64, companyID int64) (*models.Company, error)
	CreateCompanyMember(ctx context.Context, companyID int64, name, email, role string) (models.CompanyMember, error)
	CreateTeam(ctx context.Context, companyID int64, name, description string) (models.Team, error)
	GetTeamByID(ctx context.Context, companyID int64, teamID int64) (models.Team, error)
	GetTeamMembers(ctx context.Context, companyID int64, teamID int64) ([]models.CompanyMember, error)
	AddTeamMember(ctx context.Context, teamID int64, memberID int64) error
	RemoveTeamMember(ctx context.Context, companyID int64, teamID int64, memberID int64) error
	TeamExistsForCompany(ctx context.Context, teamID int64, companyID int64) (bool, error)
	CompanyMemberExistsForCompany(ctx context.Context, memberID int64, companyID int64) (bool, error)
	CreateToken(ctx context.Context, candidateID *int64, companyMemberID *int64, token string, expiresAt time.Time) error
	GetToken(ctx context.Context, token string) (*models.MCPToken, error)
	GetValidToken(ctx context.Context, token string) (*models.MCPToken, error)
	ConsumeToken(ctx context.Context, token string) (*models.MCPToken, error)
	CreateCandidate(ctx context.Context, companyID int64, name string, email string) (int64, error)
	GetCandidatesByCompany(ctx context.Context, companyID int64) ([]models.Candidate2, error)
	GetCompanyMembers(ctx context.Context, companyID int64) ([]models.CompanyMember, error)
	GetCandidates(ctx context.Context, companyID int64) ([]models.Candidate2, error)
}

var _ RepositoryStore = (*Repository)(nil)
