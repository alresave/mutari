package matching

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMatches(ctx context.Context, userID int64) ([]models.Candidate, error) {
	return s.repo.GetMCPMatches(ctx, userID)
}

func (s *Service) SubmitMCPAnswers(ctx context.Context, userID int64, answers []models.MCPAnswer) error {

	err := s.repo.mcp.StoreMCPAnswers(ctx, userID, answers)
	if err != nil {
		return err
	}

	return s.repo.mcp.BuildIdentityVector(ctx, userID)
}

func (s *Service) GetMCPQuestions(ctx context.Context) ([]models.MCPQuestion, error) {
	return s.repo.mcp.GetMCPQuestions(ctx)
}

func (s *Service) RebuildAllProfiles(ctx context.Context) error {

	ids, err := s.repo.users.GetAllUserIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {

		err := s.repo.mcp.BuildIdentityVector(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *Service) Signup(ctx context.Context, name, email, password string) (int64, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.repo.users.CreateUser(ctx, name, email, string(hash))
}
func (s *Service) Login(ctx context.Context, email, password string) (string, int64, string, error) {

	id, hash, companyID, name, err := s.repo.users.GetUserByEmail(ctx, email)
	if err != nil {
		return "", 0, "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	if err != nil {
		return "", 0, "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":    id,
		"company_id": companyID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, "", err
	}

	return tokenString, id, name, nil
}

func (s *Service) Me(
	ctx context.Context,
	userID int64,
) (int64, string, string, bool, error) {

	id, name, photo, err := s.repo.users.GetUserByID(ctx, userID)
	if err != nil {
		return 0, "", "", false, err
	}

	profileComplete, err := s.repo.mcp.HasMCPProfile(ctx, userID)
	if err != nil {
		return 0, "", "", false, err
	}

	return id, name, photo, profileComplete, nil
}

func (s *Service) GetTeamMatches(ctx context.Context, teamID int64) ([]models.Candidate, error) {
	return s.repo.GetTeamMatches(ctx, teamID)
}
func (s *Service) GetTeamFit(ctx context.Context, teamID int64, userID int64) (models.Candidate, error) {
	return s.repo.GetTeamFit(ctx, teamID, userID)
}

func (s *Service) SimulateTeamFit(
	ctx context.Context,
	teamID int64,
	candidateID int64,
) (models.TeamFitSimulation, error) {

	result, err := s.repo.SimulateTeamFit(
		ctx,
		teamID,
		candidateID,
	)

	if err != nil {
		return models.TeamFitSimulation{}, err
	}

	return result, nil
}

func (s *Service) GetTeamCulture(ctx context.Context, teamID int64) (*models.TeamCulture, error) {

	scores, archetype, err := s.repo.GetTeamCulture(ctx, teamID)
	if err != nil {
		return nil, err
	}

	result := &models.TeamCulture{
		TeamID:      teamID,
		Archetype:   archetype,
		BlockScores: scores,
	}

	return result, nil
}

func (s *Service) GetTeamGaps(ctx context.Context, teamID int64) (*models.TeamGaps, error) {

	weak, profile, err := s.repo.GetTeamGaps(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &models.TeamGaps{
		TeamID:             teamID,
		WeakBlocks:         weak,
		RecommendedProfile: profile,
	}, nil
}
func (s *Service) GetTeams(ctx context.Context) ([]models.Team, error) {
	return s.repo.GetTeams(ctx)
}

func (s *Service) RebuildAllTeamProfiles(ctx context.Context) error {

	teamIDs, err := s.repo.GetAllTeamIDs(ctx)
	if err != nil {
		return err
	}

	for _, teamID := range teamIDs {

		err := s.repo.BuildTeamVector(ctx, teamID)
		if err != nil {
			return err
		}

	}

	return nil
}
func (s *Service) AddUserToTeam(ctx context.Context, teamID int64, userID int64) error {

	err := s.repo.AddUserToTeam(ctx, teamID, userID)
	if err != nil {
		return err
	}

	return s.repo.BuildTeamVector(ctx, teamID)
}
func (s *Service) GetTeamProfile(ctx context.Context, teamID int64) ([]float32, error) {
	return s.repo.GetTeamProfile(ctx, teamID)
}
func (s *Service) GetTeamHiringProfile(ctx context.Context, teamID int64) (*models.HiringProfile, error) {

	vector, err := s.repo.GetTeamHiringProfile(ctx, teamID)
	if err != nil {
		return nil, err
	}

	blocks := make(map[string]float64)

	for i := 0; i < 8; i++ {

		start := i * 8
		end := start + 8

		sum := 0.0

		for j := start; j < end; j++ {
			sum += float64(vector[j])
		}

		avg := sum / 8.0

		blocks[models.BlockNames[i]] = avg
	}

	profile := &models.HiringProfile{
		TeamID:            teamID,
		IdealProfile:      blocks,
		SuggestedRole:     PredictRole(blocks),
		RecommendedTraits: GenerateTraits(blocks),
	}

	return profile, nil
}

// PredictRole infers the most suitable role for a candidate
// based on the averaged block scores of a team profile.

// GenerateTraits converts the psychological block profile
// into readable HR traits for the hiring profile response.

func (s *Service) GetTeamVector(ctx context.Context, teamID int64) (*models.TeamVector, error) {

	return s.repo.GetTeamVector(ctx, teamID)

}
func (s *Service) CreateCompany(ctx context.Context, name, website string, owner int64) (models.Company, error) {
	return s.repo.users.CreateCompany(ctx, name, website, owner)
}
func (s *Service) CreateCompanyMember(
	ctx context.Context,
	companyID int64,
	name, email, role string,
) (models.CompanyMember, error) {

	return s.repo.users.CreateCompanyMember(ctx, companyID, name, email, role)
}
func (s *Service) CreateTeam(
	ctx context.Context,
	companyID int64,
	name, description string,
) (models.Team, error) {

	return s.repo.users.CreateTeam(ctx, companyID, name, description)
}
func (s *Service) AddTeamMember(
	ctx context.Context,
	companyID int64,
	teamID int64,
	memberID int64,
) error {

	// 🔒 validar team pertenece a la empresa
	exists, err := s.repo.users.TeamExistsForCompany(ctx, teamID, companyID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("team not found")
	}

	// 🔒 validar member pertenece a la empresa (también importante)
	memberExists, err := s.repo.users.CompanyMemberExistsForCompany(ctx, memberID, companyID)
	if err != nil {
		return err
	}
	if !memberExists {
		return errors.New("member not found")
	}

	// 🔥 insertar
	err = s.repo.users.AddTeamMember(ctx, teamID, memberID)
	if err != nil {

		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("member already in team")
		}

		return err
	}

	return nil
}

func (s *Service) GenerateToken(ctx context.Context, candidateID int64) (string, error) {

	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	err = s.repo.users.CreateToken(ctx, candidateID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}
func (s *Service) CreateCandidate(
	ctx context.Context,
	companyID int64,
	name string,
	email string,
) (int64, error) {

	return s.repo.users.CreateCandidate(ctx, companyID, name, email)
}
