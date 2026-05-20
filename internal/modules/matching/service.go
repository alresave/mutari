package matching

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	mcprepo "github.com/alresave/mutari/internal/modules/mcp"
	"github.com/alresave/mutari/internal/modules/models"
	userrepo "github.com/alresave/mutari/internal/modules/users"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo  RepositoryStore
	users userrepo.RepositoryStore
	mcp   mcprepo.RepositoryStore
}

// NewService wires the matching service to its repository dependencies.
func NewService(repo RepositoryStore) *Service {
	return &Service{
		repo:  repo,
		users: repo.Users(),
		mcp:   repo.MCP(),
	}
}

func (s *Service) GetMatches(ctx context.Context, memberID int64) ([]models.Candidate, error) {
	return s.repo.GetMemberMatches(ctx, memberID)
}

func (s *Service) SubmitMCPAnswers(ctx context.Context, memberID int64, answers []models.MCPAnswer) error {

	err := s.mcp.StoreCompanyMemberAnswers(ctx, memberID, answers)
	if err != nil {
		return err
	}

	return s.mcp.BuildCompanyMemberProfile(ctx, memberID)
}

func (s *Service) GetMCPQuestions(ctx context.Context) ([]models.MCPQuestion, error) {
	return s.mcp.GetMCPQuestions(ctx)
}

func (s *Service) RebuildAllProfiles(ctx context.Context) error {

	ids, err := s.mcp.GetAllCompanyMemberIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {

		err := s.mcp.BuildCompanyMemberProfile(ctx, id)
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

	return s.users.CreateUser(ctx, name, email, string(hash))
}
func (s *Service) Login(ctx context.Context, email, password string) (string, int64, string, int64, error) {

	id, hash, companyID, name, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		return "", 0, "", 0, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	if err != nil {
		return "", 0, "", 0, errors.New("invalid credentials")
	}

	tokenString, err := GenerateJWTWithCompany(id, companyID)
	if err != nil {
		return "", 0, "", 0, err
	}

	return tokenString, id, name, companyID, nil
}

func (s *Service) Me(
	ctx context.Context,
	userID int64,
	companyID int64,
) (int64, string, string, *models.Company, bool, error) {

	id, name, photo, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return 0, "", "", nil, false, err
	}

	var company *models.Company

	if companyID != 0 {
		company, err = s.users.GetCompanyForUser(ctx, userID, companyID)
		if err != nil {
			return 0, "", "", nil, false, err
		}
	} else {
		company, err = s.users.GetDefaultCompanyForUser(ctx, userID)
		if err != nil {
			return 0, "", "", nil, false, err
		}
	}

	profileComplete, err := s.mcp.HasMCPProfile(ctx, userID)
	if err != nil {
		return 0, "", "", nil, false, err
	}

	return id, name, photo, company, profileComplete, nil
}

func (s *Service) GetTeamMatches(ctx context.Context, teamID int64) ([]models.Candidate, error) {
	return s.repo.GetTeamMatches(ctx, teamID)
}
func (s *Service) GetTeamMemberFit(ctx context.Context, teamID int64, memberID int64) (models.Candidate, error) {
	return s.repo.GetTeamMemberFit(ctx, teamID, memberID)
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
func (s *Service) GetTeams(ctx context.Context, companyID int64) ([]models.Team, error) {
	return s.repo.GetTeams(ctx, companyID)
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
func (s *Service) AddMemberToTeam(ctx context.Context, teamID int64, memberID int64) error {

	err := s.repo.AddMemberToTeam(ctx, teamID, memberID)
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

	if len(vector) < models.MCPVectorSize {
		return nil, fmt.Errorf("team hiring profile has %d dimensions, expected %d", len(vector), models.MCPVectorSize)
	}

	for i := 0; i < models.MCPBlockCount; i++ {

		start := i * models.MCPDimensionsPerBlock
		end := start + models.MCPDimensionsPerBlock

		sum := 0.0

		for j := start; j < end; j++ {
			sum += float64(vector[j])
		}

		avg := sum / float64(models.MCPDimensionsPerBlock)

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

func (s *Service) GetTeamVector(ctx context.Context, teamID int64) (*models.TeamVector, error) {

	return s.repo.GetTeamVector(ctx, teamID)

}
func (s *Service) CreateCompany(ctx context.Context, name, website string, owner int64) (models.Company, error) {
	return s.users.CreateCompany(ctx, name, website, owner)
}

func (s *Service) ListCompaniesForUser(ctx context.Context, userID int64) ([]models.CompanyMembership, error) {
	return s.users.ListCompaniesForUser(ctx, userID)
}

func (s *Service) SwitchCompany(ctx context.Context, userID int64, companyID int64) (*models.CompanyMembership, string, error) {
	membership, err := s.users.SetDefaultCompanyForUser(ctx, userID, companyID)
	if err != nil {
		return nil, "", err
	}

	token, err := GenerateJWTWithCompany(userID, companyID)
	if err != nil {
		return nil, "", err
	}

	return membership, token, nil
}

func (s *Service) CreateCompanyMember(
	ctx context.Context,
	companyID int64,
	name, email, role string,
) (models.CompanyMember, error) {

	return s.users.CreateCompanyMember(ctx, companyID, name, email, role)
}
func (s *Service) CreateTeam(
	ctx context.Context,
	companyID int64,
	name, description string,
) (models.Team, error) {

	return s.users.CreateTeam(ctx, companyID, name, description)
}
func (s *Service) GetTeam(ctx context.Context, companyID, teamID int64) (models.Team, error) {
	return s.users.GetTeamByID(ctx, companyID, teamID)
}

func (s *Service) GetTeamMembers(ctx context.Context, companyID, teamID int64) ([]models.CompanyMember, error) {
	return s.users.GetTeamMembers(ctx, companyID, teamID)
}

func (s *Service) RemoveTeamMember(ctx context.Context, companyID, teamID, memberID int64) error {
	err := s.users.RemoveTeamMember(ctx, companyID, teamID, memberID)
	if err != nil {
		return err
	}

	return s.repo.BuildTeamVector(ctx, teamID)
}
func (s *Service) AddTeamMember(
	ctx context.Context,
	companyID int64,
	teamID int64,
	memberID int64,
) error {

	// Guard against cross-company membership changes by checking both sides
	// against the company claim from the current JWT.
	exists, err := s.users.TeamExistsForCompany(ctx, teamID, companyID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("team not found")
	}

	memberExists, err := s.users.CompanyMemberExistsForCompany(ctx, memberID, companyID)
	if err != nil {
		return err
	}
	if !memberExists {
		return errors.New("member not found")
	}

	err = s.users.AddTeamMember(ctx, teamID, memberID)
	if err != nil {

		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("member already in team")
		}

		return err
	}

	return nil
}

func (s *Service) CreateCandidate(
	ctx context.Context,
	companyID int64,
	name string,
	email string,
) (int64, error) {

	return s.users.CreateCandidate(ctx, companyID, name, email)
}

func (s *Service) ValidateToken(
	ctx context.Context,
	token string,
) (*models.MCPToken, error) {

	if token == "" {
		return nil, errors.New("token is required")
	}

	t, err := s.users.GetValidToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return t, nil
}
func (s *Service) GenerateToken(
	ctx context.Context,
	candidateID *int64,
	companyMemberID *int64,
) (string, error) {

	// A quiz token must target exactly one entity type.
	if (candidateID == nil && companyMemberID == nil) ||
		(candidateID != nil && companyMemberID != nil) {
		return "", errors.New("invalid target")
	}

	expiresAt := time.Now().UTC().Add(48 * time.Hour)

	// UUID collisions are unlikely, but retrying keeps token creation resilient
	// when the database unique constraint rejects a generated value.
	for i := 0; i < 3; i++ {
		token := uuid.NewString()

		err := s.users.CreateToken(
			ctx,
			candidateID,
			companyMemberID,
			token,
			expiresAt,
		)

		if err == nil {
			return token, nil
		}
	}

	return "", errors.New("failed to generate token")
}
func BuildVector(answers []models.AnswerInput) []float32 {

	// The MCP model currently stores one question per vector dimension.
	vector := make([]float32, models.MCPVectorSize)

	for _, a := range answers {
		if a.QuestionID < 1 || a.QuestionID > int64(len(vector)) {
			continue
		}

		vector[a.QuestionID-1] = float32(a.Value)
	}

	return vector
}
func (s *Service) ProcessAnswers(
	ctx context.Context,
	token *models.MCPToken,
	answers []models.AnswerInput,
) error {

	// Persist answers in the table that matches the token target.
	switch token.EntityType {
	case "candidate":
		err := s.mcp.SaveAnswers(ctx, token.EntityID, toAnswers(answers))
		if err != nil {
			return err
		}
	case "member":
		err := s.mcp.SaveMemberAnswers(ctx, token.EntityID, answers)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid entity type")
	}

	// Build a 64-dimensional pgvector-compatible profile from question values.
	vector := BuildVector(answers)

	// Candidates update only their own profile; member updates also refresh team
	// vectors because team profiles are averages of member profiles.
	if token.EntityType == "candidate" {

		err := s.mcp.UpsertCandidateProfile(ctx, token.EntityID, vector)
		if err != nil {
			return err
		}

	} else {

		err := s.mcp.UpsertCompanyMemberProfile(ctx, token.EntityID, vector)
		if err != nil {
			return err
		}

		if err := s.repo.RebuildAllTeams(ctx); err != nil {
			return err
		}
	}

	// Mark the one-time quiz token as consumed after profile persistence succeeds.
	_, err := s.users.ConsumeToken(ctx, token.Token)
	if err != nil {
		return err
	} else {
		return nil
	}
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

	// Register every block from the question bank so unanswered blocks still
	// occupy a stable position in the returned vector.
	blockSet := make(map[int64]struct{})

	for _, q := range questions {

		blockSet[q.Block] = struct{}{}

		if val, ok := answerMap[q.ID]; ok {
			blockSum[q.Block] += val
			blockCount[q.Block]++
		}
	}

	// Sort block IDs to keep vector positions deterministic across map iteration.
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

func BuildBlockScores(answers []models.Answer) map[string]float64 {

	blockMap := map[int][]float64{}

	for _, a := range answers {
		blockMap[a.Block] = append(blockMap[a.Block], a.Value)
	}

	result := map[string]float64{}

	for block, values := range blockMap {

		var sum float64
		for _, v := range values {
			sum += v
		}

		avg := sum / float64(len(values)) / 5.0

		switch block {
		case 1:
			result["ownership"] = avg
		case 2:
			result["communication"] = avg
		case 3:
			result["execution"] = avg
		default:
			result[fmt.Sprintf("block_%d", block)] = avg
		}
	}

	return result
}
func BuildCandidateVector(answers []models.Answer) []float32 {

	// Store eight MCP blocks with eight dimensions per block.
	vec := make([]float32, models.MCPVectorSize)

	// Group answers by block before averaging them into block-level scores.
	blockMap := map[int][]float64{}

	for _, a := range answers {
		blockMap[a.Block] = append(blockMap[a.Block], a.Value)
	}

	// Keep block ordering stable so saved vectors remain comparable over time.
	blockOrder := make([]int, 0, models.MCPBlockCount)
	for block := 1; block <= models.MCPBlockCount; block++ {
		blockOrder = append(blockOrder, block)
	}

	index := 0

	for _, block := range blockOrder {

		values := blockMap[block]

		if len(values) == 0 {
			continue
		}

		var sum float64
		for _, v := range values {
			sum += v
		}

		avg := sum / float64(len(values))

		normalized := avg / 5.0

		// Expand each block average across its eight dimensions. This preserves
		// the 64-dimensional shape expected by team and candidate matching.
		for i := 0; i < models.MCPDimensionsPerBlock; i++ {

			if index >= len(vec) {
				break
			}

			vec[index] = float32(normalized)
			index++
		}
	}

	return vec
}

func (s *Service) SubmitAnswers(
	ctx context.Context,
	token string,
	answers []models.Answer,
) (*models.QuizResult, error) {

	// Validate the one-time token before accepting quiz answers.
	t, err := s.users.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if t.Used {
		return nil, fmt.Errorf("token already used")
	}

	if time.Now().After(t.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	candidateID := t.EntityID

	// Load the question-to-block mapping so submitted answers can be enriched
	// before they are saved and scored.
	questions, err := s.mcp.GetMCPQuestions(ctx)
	if err != nil {
		return nil, err
	}

	blockMap := make(map[int64]int)

	for _, q := range questions {
		blockMap[q.ID] = int(q.Block)
	}

	for i := range answers {

		block, ok := blockMap[answers[i].QuestionID]
		if !ok {
			return nil, fmt.Errorf("invalid question id")
		}

		answers[i].Block = block
	}

	err = s.mcp.SaveAnswers(ctx, candidateID, answers)
	if err != nil {
		return nil, err
	}

	vector := BuildCandidateVector(answers)

	blockScores := BuildBlockScores(answers)

	// Save both the profile vector used for matching and the consumed-token state
	// so the same quiz link cannot overwrite the profile later.
	err = s.mcp.UpsertCandidateProfile(ctx, candidateID, vector)
	if err != nil {
		return nil, err
	}

	if _, err = s.users.ConsumeToken(ctx, token); err != nil {
		return nil, err
	}

	return &models.QuizResult{
		CandidateID: candidateID,
		BlockScores: blockScores,
	}, nil
}
func (s *Service) GetCandidates(ctx context.Context, companyID int64) ([]models.Candidate2, error) {

	if companyID == 0 {
		return nil, fmt.Errorf("invalid company")
	}

	return s.users.GetCandidatesByCompany(ctx, companyID)
}

func (s *Service) GetCompanyMembers(ctx context.Context, companyID int64) ([]models.CompanyMember, error) {
	return s.users.GetCompanyMembers(ctx, companyID)
}

func (s *Service) GetCandidates2(ctx context.Context, companyID int64) ([]models.Candidate2, error) {
	return s.users.GetCandidates(ctx, companyID)
}

func toAnswers(inputs []models.AnswerInput) []models.Answer {
	answers := make([]models.Answer, 0, len(inputs))

	for _, input := range inputs {
		answers = append(answers, models.Answer{
			QuestionID: input.QuestionID,
			Value:      input.Value,
		})
	}

	return answers
}
