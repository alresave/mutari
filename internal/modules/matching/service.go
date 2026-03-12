package matching

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}
type TeamCulture struct {
	TeamID      int64          `json:"team_id"`
	Archetype   string         `json:"archetype"`
	BlockScores map[string]int `json:"block_scores"`
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMatches(ctx context.Context, userID int64) ([]Candidate, error) {
	return s.repo.GetMCPMatches(ctx, userID)
}

func (s *Service) SubmitMCPAnswers(ctx context.Context, userID int64, answers []MCPAnswer) error {

	err := s.repo.StoreMCPAnswers(ctx, userID, answers)
	if err != nil {
		return err
	}

	return s.repo.BuildIdentityVector(ctx, userID)
}

func (s *Service) GetMCPQuestions(ctx context.Context) ([]MCPQuestion, error) {
	return s.repo.GetMCPQuestions(ctx)
}

func (s *Service) RebuildAllProfiles(ctx context.Context) error {

	ids, err := s.repo.GetAllUserIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {

		err := s.repo.BuildIdentityVector(ctx, id)
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

	return s.repo.CreateUser(ctx, name, email, string(hash))
}
func (s *Service) Login(ctx context.Context, email, password string) (int64, string, bool, error) {

	id, hash, name, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, "", false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return 0, "", false, errors.New("invalid credentials")
	}

	hasProfile, err := s.repo.HasMCPProfile(ctx, id)
	if err != nil {
		return 0, "", false, err
	}

	return id, name, hasProfile, nil
}

var jwtSecret = []byte("mutari-secret")

func GenerateJWT(userID int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSecret)
}
func (s *Service) Me(
	ctx context.Context,
	userID int64,
) (int64, string, string, bool, error) {

	id, name, photo, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return 0, "", "", false, err
	}

	profileComplete, err := s.repo.HasMCPProfile(ctx, userID)
	if err != nil {
		return 0, "", "", false, err
	}

	return id, name, photo, profileComplete, nil
}

func (s *Service) GetTeamMatches(ctx context.Context, teamID int64) ([]Candidate, error) {
	return s.repo.GetTeamMatches(ctx, teamID)
}
func (s *Service) GetTeamFit(ctx context.Context, teamID int64, userID int64) (Candidate, error) {
	return s.repo.GetTeamFit(ctx, teamID, userID)
}

func (s *Service) SimulateTeamFit(
	ctx context.Context,
	teamID int64,
	candidateID int64,
) (TeamFitSimulation, error) {

	result, err := s.repo.SimulateTeamFit(
		ctx,
		teamID,
		candidateID,
	)

	if err != nil {
		return TeamFitSimulation{}, err
	}

	return result, nil
}

func (s *Service) GetTeamCulture(ctx context.Context, teamID int64) (*TeamCulture, error) {

	scores, archetype, err := s.repo.GetTeamCulture(ctx, teamID)
	if err != nil {
		return nil, err
	}

	result := &TeamCulture{
		TeamID:      teamID,
		Archetype:   archetype,
		BlockScores: scores,
	}

	return result, nil
}

// team average vector
