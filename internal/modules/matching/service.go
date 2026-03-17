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
type HiringProfile struct {
	TeamID            int64              `json:"team_id"`
	IdealProfile      map[string]float64 `json:"ideal_profile"`
	SuggestedRole     string             `json:"suggested_role"`
	RecommendedTraits []string           `json:"recommended_traits"`
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
func (s *Service) Login(ctx context.Context, email, password string) (string, int64, string, error) {

	id, hash, companyID, name, err := s.repo.GetUserByEmail(ctx, email)
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

type TeamGaps struct {
	TeamID             int64    `json:"team_id"`
	WeakBlocks         []string `json:"weak_blocks"`
	RecommendedProfile string   `json:"recommended_profile"`
}

func (s *Service) GetTeamGaps(ctx context.Context, teamID int64) (*TeamGaps, error) {

	weak, profile, err := s.repo.GetTeamGaps(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &TeamGaps{
		TeamID:             teamID,
		WeakBlocks:         weak,
		RecommendedProfile: profile,
	}, nil
}
func (s *Service) GetTeams(ctx context.Context) ([]Team, error) {
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
func (s *Service) GetTeamHiringProfile(ctx context.Context, teamID int64) (*HiringProfile, error) {

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

		blocks[blockNames[i]] = avg
	}

	profile := &HiringProfile{
		TeamID:            teamID,
		IdealProfile:      blocks,
		SuggestedRole:     PredictRole(blocks),
		RecommendedTraits: GenerateTraits(blocks),
	}

	return profile, nil
}

// PredictRole infers the most suitable role for a candidate
// based on the averaged block scores of a team profile.
func PredictRole(blocks map[string]float64) string {

	ambition := blocks["Ambition"]
	cognitive := blocks["Cognitive Style"]
	emotional := blocks["Emotional Stability"]
	relational := blocks["Relational Architecture"]
	growth := blocks["Growth"]

	// Simple heuristics that work well for demos and HR explanations

	if ambition > 0.85 && emotional > 0.80 {
		return "Execution Leader"
	}

	if cognitive > 0.88 && growth > 0.80 {
		return "Strategic Architect"
	}

	if relational > 0.85 && emotional > 0.80 {
		return "Team Integrator"
	}

	if cognitive > 0.90 {
		return "Specialist"
	}

	return "Balanced Contributor"
}

// GenerateTraits converts the psychological block profile
// into readable HR traits for the hiring profile response.
func GenerateTraits(blocks map[string]float64) []string {

	traits := []string{}

	if blocks["Cognitive Style"] > 0.85 {
		traits = append(traits, "Strategic thinker")
	}

	if blocks["Emotional Stability"] > 0.80 {
		traits = append(traits, "High emotional regulation")
	}

	if blocks["Ambition"] > 0.85 {
		traits = append(traits, "Strong execution drive")
	}

	if blocks["Relational Architecture"] > 0.80 {
		traits = append(traits, "Collaborative mindset")
	}

	if blocks["Growth"] > 0.80 {
		traits = append(traits, "Continuous improvement oriented")
	}

	if blocks["Core Values"] > 0.80 {
		traits = append(traits, "Values aligned with team culture")
	}

	// fallback in case no traits triggered
	if len(traits) == 0 {
		traits = append(traits, "Balanced team contributor")
	}

	return traits
}
