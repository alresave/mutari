package matching

import (
	"context"

	mcprepo "github.com/alresave/mutari/internal/modules/mcp"
	"github.com/alresave/mutari/internal/modules/models"
	userrepo "github.com/alresave/mutari/internal/modules/users"
)

// RepositoryStore is the matching service's persistence boundary.
type RepositoryStore interface {
	Users() userrepo.RepositoryStore
	MCP() mcprepo.RepositoryStore
	GetMemberMatches(ctx context.Context, memberID int64) ([]models.Candidate, error)
	GetTeamMatches(ctx context.Context, teamID int64) ([]models.Candidate, error)
	GetTeamMemberFit(ctx context.Context, teamID int64, memberID int64) (models.Candidate, error)
	SimulateTeamFit(ctx context.Context, teamID int64, candidateID int64) (models.TeamFitSimulation, error)
	GetTeamCulture(ctx context.Context, teamID int64) (map[string]int, string, error)
	GetTeamGaps(ctx context.Context, teamID int64) ([]string, string, error)
	GetTeams(ctx context.Context, companyID int64) ([]models.Team, error)
	GetAllTeamIDs(ctx context.Context) ([]int64, error)
	BuildTeamVector(ctx context.Context, teamID int64) error
	AddMemberToTeam(ctx context.Context, teamID int64, memberID int64) error
	GetTeamProfile(ctx context.Context, teamID int64) ([]float32, error)
	GetTeamHiringProfile(ctx context.Context, teamID int64) ([]float32, error)
	GetTeamVector(ctx context.Context, teamID int64) (*models.TeamVector, error)
	RebuildAllTeams(ctx context.Context) error
}

// ServiceStore is the HTTP handler's service boundary.
type ServiceStore interface {
	GetMatches(ctx context.Context, memberID int64) ([]models.Candidate, error)
	SubmitMCPAnswers(ctx context.Context, memberID int64, answers []models.MCPAnswer) error
	GetMCPQuestions(ctx context.Context) ([]models.MCPQuestion, error)
	RebuildAllProfiles(ctx context.Context) error
	Signup(ctx context.Context, name, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, int64, string, int64, error)
	Me(ctx context.Context, userID int64, companyID int64) (int64, string, string, *models.Company, bool, error)
	GetTeamMatches(ctx context.Context, teamID int64) ([]models.Candidate, error)
	GetTeamMemberFit(ctx context.Context, teamID int64, memberID int64) (models.Candidate, error)
	SimulateTeamFit(ctx context.Context, teamID int64, candidateID int64) (models.TeamFitSimulation, error)
	GetTeamCulture(ctx context.Context, teamID int64) (*models.TeamCulture, error)
	GetTeamGaps(ctx context.Context, teamID int64) (*models.TeamGaps, error)
	GetTeams(ctx context.Context, companyID int64) ([]models.Team, error)
	RebuildAllTeamProfiles(ctx context.Context) error
	AddMemberToTeam(ctx context.Context, teamID int64, memberID int64) error
	GetTeamProfile(ctx context.Context, teamID int64) ([]float32, error)
	GetTeamHiringProfile(ctx context.Context, teamID int64) (*models.HiringProfile, error)
	GetTeamVector(ctx context.Context, teamID int64) (*models.TeamVector, error)
	CreateCompany(ctx context.Context, name, website string, owner int64) (models.Company, error)
	ListCompaniesForUser(ctx context.Context, userID int64) ([]models.CompanyMembership, error)
	SwitchCompany(ctx context.Context, userID int64, companyID int64) (*models.CompanyMembership, string, error)
	CreateCompanyMember(ctx context.Context, companyID int64, name, email, role string) (models.CompanyMember, error)
	CreateTeam(ctx context.Context, companyID int64, name, description string) (models.Team, error)
	GetTeam(ctx context.Context, companyID, teamID int64) (models.Team, error)
	GetTeamMembers(ctx context.Context, companyID, teamID int64) ([]models.CompanyMember, error)
	RemoveTeamMember(ctx context.Context, companyID, teamID, memberID int64) error
	AddTeamMember(ctx context.Context, companyID int64, teamID int64, memberID int64) error
	CreateCandidate(ctx context.Context, companyID int64, name string, email string) (int64, error)
	ValidateToken(ctx context.Context, token string) (*models.MCPToken, error)
	GenerateToken(ctx context.Context, candidateID *int64, companyMemberID *int64) (string, error)
	ProcessAnswers(ctx context.Context, token *models.MCPToken, answers []models.AnswerInput) error
	SubmitAnswers(ctx context.Context, token string, answers []models.Answer) (*models.QuizResult, error)
	GetCandidates(ctx context.Context, companyID int64) ([]models.Candidate2, error)
	GetCompanyMembers(ctx context.Context, companyID int64) ([]models.CompanyMember, error)
	GetCandidates2(ctx context.Context, companyID int64) ([]models.Candidate2, error)
}

var _ RepositoryStore = (*Repository)(nil)
var _ ServiceStore = (*Service)(nil)
