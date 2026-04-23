package matching

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alresave/mutari/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

func testDatabaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}

	return "postgres://postgres:mutari123@localhost:5432/mutari?sslmode=disable"
}

func newTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	config, err := pgxpool.ParseConfig(testDatabaseURL())
	if err != nil {
		t.Fatalf("parse database config: %v", err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Skipf("skip integration test, database unavailable: %v", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		t.Skipf("skip integration test, database unavailable: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func resetIntegrationDB(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	query := `
	TRUNCATE TABLE
		mcp_tokens,
		candidate_mcp_answers,
		candidate_mcp_profiles,
		company_member_mcp_answers,
		company_member_mcp_profiles,
		team_members,
		team_profiles,
		candidates,
		company_members,
		teams,
		companies,
		users
	RESTART IDENTITY CASCADE
	`

	if _, err := db.Exec(context.Background(), query); err != nil {
		t.Fatalf("reset database: %v", err)
	}

	ensureQuestionBank(t, db)
}

func ensureQuestionBank(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	var count int
	if err := db.QueryRow(context.Background(), `SELECT COUNT(*) FROM mcp_questions`).Scan(&count); err != nil {
		t.Fatalf("count mcp questions: %v", err)
	}

	if count > 0 {
		return
	}

	blocks := []string{
		"Execution Drive",
		"Cognitive Processing",
		"Professional Values",
		"Emotional Resilience",
		"Learning Velocity",
		"Collaboration Depth",
		"Operational Rhythm",
		"Team Positioning",
	}

	for i, name := range blocks {
		if _, err := db.Exec(
			context.Background(),
			`INSERT INTO mcp_blocks (id, name, description) VALUES ($1, $2, $3)`,
			i+1,
			name,
			name+" description",
		); err != nil {
			t.Fatalf("insert mcp block %d: %v", i+1, err)
		}
	}

	for block := 1; block <= 8; block++ {
		for question := 1; question <= 8; question++ {
			text := fmt.Sprintf("Block %d question %d", block, question)
			if _, err := db.Exec(
				context.Background(),
				`INSERT INTO mcp_questions (block, question_text, scale) VALUES ($1, $2, 'likert_1_5')`,
				block,
				text,
			); err != nil {
				t.Fatalf("insert mcp question %d/%d: %v", block, question, err)
			}
		}
	}
}

func testAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenString := header[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", int64(userID))
		}
		if companyID, ok := claims["company_id"].(float64); ok {
			c.Set("company_id", int64(companyID))
		}

		c.Next()
	}
}

func newTestRouter(t *testing.T, db *pgxpool.Pool) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, security.NewRateLimiter())

	router := gin.New()
	router.POST("/auth/signup", handler.Signup)
	router.POST("/auth/login", handler.Login)
	router.GET("/quiz", handler.GetQuiz)

	auth := router.Group("/")
	auth.Use(testAuthMiddleware())
	auth.POST("/companies", handler.CreateCompany)
	auth.GET("/me", handler.Me)
	auth.GET("/company-members", handler.GetCompanyMembers)
	auth.POST("/company-members", handler.CreateCompanyMember)
	auth.GET("/teams", handler.GetTeams)
	auth.POST("/teams", handler.CreateTeam)
	auth.GET("/teams/:team_id/members", handler.GetTeamMembers)
	auth.POST("/teams/:team_id/members", handler.AddMemberToTeam)
	auth.GET("/teams/:team_id/matches", handler.GetTeamMatches)
	auth.GET("/teams/:team_id/gaps", handler.GetTeamGaps)
	auth.GET("/teams/:team_id/hiring-profile", handler.GetTeamHiringProfile)
	auth.POST("/tokens", handler.GenerateToken)
	auth.POST("/mcp/answers", handler.SubmitAnswers)
	auth.POST("/candidates", handler.CreateCandidate)
	auth.GET("/candidates", handler.GetCandidates2)

	return router
}

func performJSONRequest(t *testing.T, router http.Handler, method, path string, body any, token string) map[string]any {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("%s %s returned %d: %s", method, path, rec.Code, rec.Body.String())
	}

	if rec.Body.Len() == 0 {
		return map[string]any{}
	}

	var response map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response %s %s: %v", method, path, err)
	}

	return response
}

func floatID(value any) int64 {
	return int64(value.(float64))
}

func createAnsweredTeamFixture(t *testing.T, router http.Handler) (string, int64, int64) {
	t.Helper()

	signup := performJSONRequest(t, router, http.MethodPost, "/auth/signup", map[string]any{
		"name":     "Flow Owner",
		"email":    "owner@test.local",
		"password": "Secret123!",
	}, "")

	userToken := signup["token"].(string)

	company := performJSONRequest(t, router, http.MethodPost, "/companies", map[string]any{
		"name":    "Flow Co",
		"website": "https://flow.example.com",
	}, userToken)

	companyToken := company["token"].(string)

	member1 := performJSONRequest(t, router, http.MethodPost, "/company-members", map[string]any{
		"name":  "Member One",
		"email": "member1@test.local",
		"role":  "Engineer",
	}, companyToken)
	member2 := performJSONRequest(t, router, http.MethodPost, "/company-members", map[string]any{
		"name":  "Member Two",
		"email": "member2@test.local",
		"role":  "Designer",
	}, companyToken)

	team := performJSONRequest(t, router, http.MethodPost, "/teams", map[string]any{
		"name":        "Team One",
		"description": "Integration test team",
	}, companyToken)

	teamID := floatID(team["id"])
	member1ID := floatID(member1["id"])
	member2ID := floatID(member2["id"])

	performJSONRequest(t, router, http.MethodPost, fmt.Sprintf("/teams/%d/members", teamID), map[string]any{
		"member_id": member1ID,
	}, companyToken)
	performJSONRequest(t, router, http.MethodPost, fmt.Sprintf("/teams/%d/members", teamID), map[string]any{
		"member_id": member2ID,
	}, companyToken)

	for _, memberID := range []int64{member1ID, member2ID} {
		tokenResp := performJSONRequest(t, router, http.MethodPost, "/tokens", map[string]any{
			"company_member_id": memberID,
		}, companyToken)

		quizResp := performJSONRequest(t, router, http.MethodGet, "/quiz?token="+tokenResp["token"].(string), nil, "")
		questions := quizResp["questions"].([]any)

		answers := make([]map[string]any, 0, len(questions))
		for _, rawQuestion := range questions {
			question := rawQuestion.(map[string]any)
			answers = append(answers, map[string]any{
				"question_id": floatID(question["id"]),
				"value":       3,
			})
		}

		performJSONRequest(t, router, http.MethodPost, "/mcp/answers", map[string]any{
			"token":   tokenResp["token"].(string),
			"answers": answers,
		}, companyToken)
	}

	candidate := performJSONRequest(t, router, http.MethodPost, "/candidates", map[string]any{
		"name":  "Candidate One",
		"email": "candidate@test.local",
	}, companyToken)

	candidateID := floatID(candidate["id"])
	tokenResp := performJSONRequest(t, router, http.MethodPost, "/tokens", map[string]any{
		"candidate_id": candidateID,
	}, companyToken)

	quizResp := performJSONRequest(t, router, http.MethodGet, "/quiz?token="+tokenResp["token"].(string), nil, "")
	questions := quizResp["questions"].([]any)

	answers := make([]map[string]any, 0, len(questions))
	for _, rawQuestion := range questions {
		question := rawQuestion.(map[string]any)
		answers = append(answers, map[string]any{
			"question_id": floatID(question["id"]),
			"value":       4,
		})
	}

	performJSONRequest(t, router, http.MethodPost, "/mcp/answers", map[string]any{
		"token":   tokenResp["token"].(string),
		"answers": answers,
	}, companyToken)

	return companyToken, teamID, candidateID
}

func TestAuthAndCompanyFlow(t *testing.T) {
	db := newTestDB(t)
	resetIntegrationDB(t, db)

	router := newTestRouter(t, db)

	signup := performJSONRequest(t, router, http.MethodPost, "/auth/signup", map[string]any{
		"name":     "Alice Owner",
		"email":    "alice.owner@test.local",
		"password": "Secret123!",
	}, "")

	login := performJSONRequest(t, router, http.MethodPost, "/auth/login", map[string]any{
		"email":    "alice.owner@test.local",
		"password": "Secret123!",
	}, "")

	company := performJSONRequest(t, router, http.MethodPost, "/companies", map[string]any{
		"name":    "Acme Labs",
		"website": "https://acme.example.com",
	}, signup["token"].(string))

	companyToken := company["token"].(string)
	member := performJSONRequest(t, router, http.MethodPost, "/company-members", map[string]any{
		"name":  "Bob Builder",
		"email": "bob.builder@test.local",
		"role":  "Engineer",
	}, companyToken)

	team := performJSONRequest(t, router, http.MethodPost, "/teams", map[string]any{
		"name":        "Core Team",
		"description": "Main delivery team",
	}, companyToken)

	performJSONRequest(t, router, http.MethodPost, fmt.Sprintf("/teams/%d/members", floatID(team["id"])), map[string]any{
		"member_id": floatID(member["id"]),
	}, companyToken)

	teamsResp := performJSONRequest(t, router, http.MethodGet, "/teams", nil, companyToken)
	membersResp := performJSONRequest(t, router, http.MethodGet, "/company-members", nil, companyToken)

	if floatID(login["user"].(map[string]any)["id"]) != floatID(signup["user_id"]) {
		t.Fatalf("login user id does not match signup user id")
	}

	if len(teamsResp["teams"].([]any)) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teamsResp["teams"].([]any)))
	}

	if len(membersResp["members"].([]any)) != 1 {
		t.Fatalf("expected 1 company member, got %d", len(membersResp["members"].([]any)))
	}
}

func TestSubmitAnswersRebuildsTeamProfile(t *testing.T) {
	db := newTestDB(t)
	resetIntegrationDB(t, db)

	router := newTestRouter(t, db)
	companyToken, teamID, _ := createAnsweredTeamFixture(t, router)

	hiringProfile := performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/teams/%d/hiring-profile", teamID), nil, companyToken)

	var profileCount int
	if err := db.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM team_profiles WHERE team_id = $1`,
		teamID,
	).Scan(&profileCount); err != nil {
		t.Fatalf("count team profiles: %v", err)
	}

	if profileCount != 1 {
		t.Fatalf("expected team profile to be rebuilt, got %d rows", profileCount)
	}

	idealProfile := hiringProfile["ideal_profile"].(map[string]any)
	if len(idealProfile) != 8 {
		t.Fatalf("expected 8 hiring profile blocks, got %d", len(idealProfile))
	}
}

func TestDashboardEndpointsReturnData(t *testing.T) {
	db := newTestDB(t)
	resetIntegrationDB(t, db)

	router := newTestRouter(t, db)
	companyToken, teamID, _ := createAnsweredTeamFixture(t, router)

	matches := performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/teams/%d/matches", teamID), nil, companyToken)
	gaps := performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/teams/%d/gaps", teamID), nil, companyToken)
	hiringProfile := performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/teams/%d/hiring-profile", teamID), nil, companyToken)

	if len(matches["matches"].([]any)) == 0 {
		t.Fatalf("expected at least one team match candidate")
	}

	if len(gaps["weak_blocks"].([]any)) == 0 {
		t.Fatalf("expected at least one weak block in team gaps")
	}

	if hiringProfile["suggested_role"] == "" {
		t.Fatalf("expected suggested role in hiring profile")
	}
}
