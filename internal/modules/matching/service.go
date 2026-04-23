package matching

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMatches(ctx context.Context, memberID int64) ([]models.Candidate, error) {
	return s.repo.GetMemberMatches(ctx, memberID)
}

func (s *Service) SubmitMCPAnswers(ctx context.Context, memberID int64, answers []models.MCPAnswer) error {

	err := s.repo.mcp.StoreCompanyMemberAnswers(ctx, memberID, answers)
	if err != nil {
		return err
	}

	return s.repo.mcp.BuildCompanyMemberProfile(ctx, memberID)
}

func (s *Service) GetMCPQuestions(ctx context.Context) ([]models.MCPQuestion, error) {
	return s.repo.mcp.GetMCPQuestions(ctx)
}

func (s *Service) RebuildAllProfiles(ctx context.Context) error {

	ids, err := s.repo.mcp.GetAllCompanyMemberIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {

		err := s.repo.mcp.BuildCompanyMemberProfile(ctx, id)
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
func (s *Service) Login(ctx context.Context, email, password string) (string, int64, string, error, int64) {

	id, hash, companyID, name, err := s.repo.users.GetUserByEmail(ctx, email)
	if err != nil {
		return "", 0, "", err, 0
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	if err != nil {
		return "", 0, "", errors.New("invalid credentials"), 0
	}

	claims := jwt.MapClaims{
		"user_id":    id,
		"company_id": companyID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, "", err, companyID
	}

	return tokenString, id, name, nil, companyID
}

func (s *Service) Me(
	ctx context.Context,
	userID int64,
) (int64, string, string, *models.Company, bool, error) {

	id, name, photo, err := s.repo.users.GetUserByID(ctx, userID)
	if err != nil {
		return 0, "", "", nil, false, err
	}

	company, err := s.repo.users.GetCompanyForOwner(ctx, userID)
	if err != nil {
		return 0, "", "", nil, false, err
	}

	profileComplete, err := s.repo.mcp.HasMCPProfile(ctx, userID)
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
func (s *Service) GetTeam(ctx context.Context, companyID, teamID int64) (models.Team, error) {
	return s.repo.users.GetTeamByID(ctx, companyID, teamID)
}

func (s *Service) GetTeamMembers(ctx context.Context, companyID, teamID int64) ([]models.CompanyMember, error) {
	return s.repo.users.GetTeamMembers(ctx, companyID, teamID)
}

func (s *Service) RemoveTeamMember(ctx context.Context, companyID, teamID, memberID int64) error {
	err := s.repo.users.RemoveTeamMember(ctx, companyID, teamID, memberID)
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

func (s *Service) CreateCandidate(
	ctx context.Context,
	companyID int64,
	name string,
	email string,
) (int64, error) {

	return s.repo.users.CreateCandidate(ctx, companyID, name, email)
}

func (s *Service) ValidateToken(
	ctx context.Context,
	token string,
) (*models.MCPToken, error) {

	if token == "" {
		return nil, errors.New("token is required")
	}

	t, err := s.repo.users.GetValidToken(ctx, token)
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

	// Validación clara
	if (candidateID == nil && companyMemberID == nil) ||
		(candidateID != nil && companyMemberID != nil) {
		return "", errors.New("invalid target")
	}

	expiresAt := time.Now().UTC().Add(48 * time.Hour)

	// Retry por colisión (muy raro, pero correcto)
	for i := 0; i < 3; i++ {
		token := uuid.NewString()

		err := s.repo.users.CreateToken(
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

	vector := make([]float32, len(answers))

	for i, a := range answers {
		vector[i] = float32(a.Value)
	}

	return vector
}
func (s *Service) ProcessAnswers(
	ctx context.Context,
	token *models.MCPToken,
	answers []models.AnswerInput,
) error {

	// 🔹 guardar respuestas según tipo
	if token.EntityType == "candidate" {

		err := s.repo.mcp.SaveAnswers(ctx, token.EntityID, toAnswers(answers))
		if err != nil {
			return err
		}

	} else if token.EntityType == "member" {

		err := s.repo.mcp.SaveMemberAnswers(ctx, token.EntityID, answers)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("invalid entity type")
	}

	// 🔹 construir vector
	vector := BuildVector(answers)

	// 🔹 guardar profile
	if token.EntityType == "candidate" {

		err := s.repo.mcp.UpsertCandidateProfile(ctx, token.EntityID, vector)
		if err != nil {
			return err
		}

	} else {

		err := s.repo.mcp.UpsertCompanyMemberProfile(ctx, token.EntityID, vector)
		if err != nil {
			return err
		}
	}

	// 🔹 marcar token usado
	_, err := s.repo.users.ConsumeToken(ctx, token.Token)
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

	// vector final (64 dims para pgvector)
	vec := make([]float32, 64)

	// agrupar por bloque
	blockMap := map[int][]float64{}

	for _, a := range answers {
		blockMap[a.Block] = append(blockMap[a.Block], a.Value)
	}

	// orden de bloques (IMPORTANTE: consistente siempre)
	blockOrder := []int{1, 2, 3, 4, 5, 6, 7, 8}

	index := 0

	for _, block := range blockOrder {

		values := blockMap[block]

		if len(values) == 0 {
			continue
		}

		// promedio del bloque
		var sum float64
		for _, v := range values {
			sum += v
		}

		avg := sum / float64(len(values)) // escala 1–5

		// normalizar 0–1
		normalized := avg / 5.0

		// 🔥 expandir en múltiples dimensiones (más robusto)
		for i := 0; i < 8; i++ {

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

	// -----------------------------
	// 1. VALIDAR TOKEN
	// -----------------------------
	t, err := s.repo.users.GetToken(ctx, token)
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

	// -----------------------------
	// 2. MAPEAR QUESTIONS → BLOCK
	// -----------------------------
	questions, err := s.repo.mcp.GetMCPQuestions(ctx)
	if err != nil {
		return nil, err
	}

	blockMap := make(map[int64]int)

	for _, q := range questions {
		blockMap[q.ID] = int(q.Block)
	}

	// -----------------------------
	// 3. ENRIQUECER ANSWERS
	// -----------------------------
	for i := range answers {

		block, ok := blockMap[answers[i].QuestionID]
		if !ok {
			return nil, fmt.Errorf("invalid question id")
		}

		answers[i].Block = block
	}

	// -----------------------------
	// 4. GUARDAR RESPUESTAS
	// -----------------------------
	err = s.repo.mcp.SaveAnswers(ctx, candidateID, answers)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 5. BUILD VECTOR
	// -----------------------------
	vector := BuildCandidateVector(answers)

	// -----------------------------
	// 6. BLOCK SCORES
	// -----------------------------
	blockScores := BuildBlockScores(answers)

	// -----------------------------
	// 7. GUARDAR PERFIL
	// -----------------------------
	err = s.repo.mcp.UpsertCandidateProfile(ctx, candidateID, vector)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 8. MARCAR TOKEN USADO
	// -----------------------------
	_, err = s.repo.users.ConsumeToken(ctx, token)

	// -----------------------------
	// 9. RETORNAR RESULTADO
	// -----------------------------
	return &models.QuizResult{
		CandidateID: candidateID,
		BlockScores: blockScores,
	}, nil
}
func (s *Service) GetCandidates(ctx context.Context, companyID int64) ([]models.Candidate2, error) {

	if companyID == 0 {
		return nil, fmt.Errorf("invalid company")
	}

	return s.repo.users.GetCandidatesByCompany(ctx, companyID)
}

func (s *Service) GetCompanyMembers(ctx context.Context, companyID int64) ([]models.CompanyMember, error) {
	return s.repo.users.GetCompanyMembers(ctx, companyID)
}

func (s *Service) GetCandidates2(ctx context.Context, companyID int64) ([]models.Candidate2, error) {
	return s.repo.users.GetCandidates(ctx, companyID)
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
