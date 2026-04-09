package matching

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/alresave/mutari/internal/security"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      *Service
	loginLimiter *security.RateLimiter
}

func NewHandler(service *Service, limiter *security.RateLimiter) *Handler {
	return &Handler{
		service:      service,
		loginLimiter: limiter,
	}
}

//
// SIGNUP
//
// @Summary User signup
// @Description Create a new user account
// @Accept json
// @Produce json
// @Param request body SignupRequest true "Signup request"
// @Success 200 {object} map[string]interface{} "user_id and token"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Failure 500 {object} map[string]interface{} "server error"
// @Router /auth/signup [post]

func (h *Handler) Signup(c *gin.Context) {

	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.service.Signup(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, err := GenerateJWT(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	c.JSON(200, gin.H{
		"user_id": id,
		"token":   token,
	})
}

//
// LOGIN
//
// @Summary User login
// @Description Authenticate user with email and password
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} map[string]interface{} "user_id and token"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Failure 429 {object} map[string]interface{} "rate limit exceeded"
// @Failure 500 {object} map[string]interface{} "server error"
// @Router /auth/login [post]

func (h *Handler) Login(c *gin.Context) {

	var body models.LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	token, userID, name, err := h.service.Login(
		c.Request.Context(),
		body.Email,
		body.Password,
	)

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":   userID,
			"name": name,
		},
	})
}

//
// GET MCP QUESTIONS
//
// @Summary Get MCP questions
// @Description Retrieve all matching compatibility profile (MCP) questions
// @Produce json
// @Success 200 {object} map[string]interface{} "array of questions"
// @Failure 500 {object} map[string]interface{} "server error"
// @Router /mcp/questions [get]

func (h *Handler) GetMCPQuestions(c *gin.Context) {

	questions, err := h.service.GetMCPQuestions(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"questions": questions,
	})
}

//
// SUBMIT MCP ANSWERS
//

//
// GET MATCHES
//

func (h *Handler) GetMatches(c *gin.Context) {

	userID := c.GetInt64("user_id")

	matches, err := h.service.GetMatches(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user_id": userID,
		"matches": matches,
	})
}

//
// ADMIN: rebuild MCP profiles
//

func (h *Handler) RebuildProfiles(c *gin.Context) {

	err := h.service.RebuildAllProfiles(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "profiles rebuilt",
	})
}
func (h *Handler) Me(c *gin.Context) {

	userID := c.GetInt64("user_id")

	id, name, photo, company, profileComplete, err :=
		h.service.Me(c.Request.Context(), userID)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{

		"user": gin.H{
			"id":    id,
			"name":  name,
			"photo": photo,
		},

		"company":          company,
		"profile_complete": profileComplete,
	})
}

func (h *Handler) GetTeamMatches(c *gin.Context) {

	teamIDParam := c.Param("team_id")

	teamID, err := strconv.ParseInt(teamIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team_id"})
		return
	}

	matches, err := h.service.GetTeamMatches(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"team_id": teamID,
		"matches": matches,
	})
}

func (h *Handler) GetTeamFit(c *gin.Context) {

	teamParam := c.Param("team_id")
	userParam := c.Param("user_id")

	teamID, _ := strconv.ParseInt(teamParam, 10, 64)
	userID, _ := strconv.ParseInt(userParam, 10, 64)

	match, err := h.service.GetTeamFit(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, match)

}
func (h *Handler) SimulateTeamFit(c *gin.Context) {

	teamIDParam := c.Param("team_id")
	candidateParam := c.Param("candidate_id")

	teamID, err := strconv.ParseInt(teamIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid team_id",
		})
		return
	}

	candidateID, err := strconv.ParseInt(candidateParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid candidate_id",
		})
		return
	}

	result, err := h.service.SimulateTeamFit(
		c.Request.Context(),
		teamID,
		candidateID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
func (h *Handler) GetTeamCulture(c *gin.Context) {

	teamIDParam := c.Param("team_id")

	teamID, err := strconv.ParseInt(teamIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid team_id",
		})
		return
	}

	result, err := h.service.GetTeamCulture(c.Request.Context(), teamID)
	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, result)
}

func (h *Handler) GetTeamGaps(c *gin.Context) {

	param := c.Param("team_id")

	teamID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team_id"})
		return
	}

	result, err := h.service.GetTeamGaps(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
func (h *Handler) GetTeams(c *gin.Context) {

	teams, err := h.service.GetTeams(c.Request.Context())
	if err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return

	}

	c.JSON(200, gin.H{
		"teams": teams,
	})

}
func (h *Handler) RebuildTeamProfiles(c *gin.Context) {

	err := h.service.RebuildAllTeamProfiles(c.Request.Context())
	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "team profiles rebuilt",
	})
}
func (h *Handler) AddUserToTeam(c *gin.Context) {

	teamParam := c.Param("team_id")

	teamID, err := strconv.ParseInt(teamParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team id"})
		return
	}

	var body struct {
		UserID int64 `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	err = h.service.AddUserToTeam(c.Request.Context(), teamID, body.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "user added to team",
	})
}
func (h *Handler) GetTeamProfile(c *gin.Context) {

	idParam := c.Param("team_id")

	teamID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team id"})
		return
	}

	vector, err := h.service.GetTeamProfile(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"team_id": teamID,
		"vector":  vector,
	})
}
func (h *Handler) GetTeamHiringProfile(c *gin.Context) {

	teamID, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team_id"})
		return
	}

	profile, err := h.service.GetTeamHiringProfile(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, profile)
}
func (h *Handler) GetTeamVector(c *gin.Context) {

	teamIDParam := c.Param("team_id")
	log.Println("**************Team VECTOR*******************************************************")

	teamID, err := strconv.ParseInt(teamIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team id"})
		return
	}

	result, err := h.service.GetTeamVector(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
func (h *Handler) CreateCompany(c *gin.Context) {

	var req models.CreateCompanyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetInt64("user_id") // 👈 viene del JWT middleware

	company, err := h.service.CreateCompany(
		c.Request.Context(),
		req.Name,
		req.Website,
		userID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, company)
}
func (h *Handler) CreateCompanyMember(c *gin.Context) {

	var req models.CreateCompanyMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	companyID := c.GetInt64("company_id")

	member, err := h.service.CreateCompanyMember(
		c.Request.Context(),
		companyID,
		req.Name,
		req.Email,
		req.Role,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, member)
}
func (h *Handler) CreateTeam(c *gin.Context) {

	var req models.CreateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	companyID := c.GetInt64("company_id")

	team, err := h.service.CreateTeam(
		c.Request.Context(),
		companyID,
		req.Name,
		req.Description,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, team)
}

func (h *Handler) AddTeamMember(c *gin.Context) {

	// 🔹 team_id desde URL
	teamIDParam := c.Param("team_id")

	teamID, err := strconv.ParseInt(teamIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid team id"})
		return
	}

	// 🔹 body
	var req models.AddTeamMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// 🔹 company_id desde JWT (middleware)
	companyID := c.GetInt64("company_id")

	// 🔹 lógica
	err = h.service.AddTeamMember(
		c.Request.Context(),
		companyID,
		teamID,
		req.CompanyMemberID,
	)

	if err != nil {

		switch err.Error() {
		case "team not found":
			c.JSON(404, gin.H{"error": errors.New("team not found " + teamIDParam + strconv.FormatInt(companyID, 10)).Error()})
			return
		case "member not found":
			c.JSON(404, gin.H{"error": err.Error()})
			return
		case "member already in team":
			c.JSON(400, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// ✅ success
	c.JSON(200, gin.H{
		"status": "member added",
	})
}

func (h *Handler) GenerateToken(c *gin.Context) {

	var req models.GenerateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.service.GenerateToken(
		c.Request.Context(),
		req.CandidateID,
		req.CompanyMemberID,
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"url":   fmt.Sprintf("http://localhost:8080/quiz.html?token=%s", token),
	})
}

func (h *Handler) CreateCandidate(c *gin.Context) {

	var req models.CreateCandidateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	companyID := c.GetInt64("company_id")

	id, err := h.service.CreateCandidate(
		c.Request.Context(),
		companyID,
		req.Name,
		req.Email,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}

func (h *Handler) GetQuiz(c *gin.Context) {

	token := c.Query("token")

	if token == "" {
		c.JSON(400, gin.H{"error": "missing token"})
		return
	}

	// 🔹 validar token
	t, err := h.service.ValidateToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔹 obtener preguntas MCP
	questions, err := h.service.GetMCPQuestions(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load questions"})
		return
	}

	// 🔥 respuesta correcta (nuevo modelo)
	c.JSON(200, gin.H{
		"entity_type": t.EntityType,
		"entity_id":   t.EntityID,
		"questions":   questions,
	})
}
func (h *Handler) SubmitAnswers(c *gin.Context) {

	var req models.SubmitAnswersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// 🔹 validar token
	t, err := h.service.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔹 procesar respuestas
	err = h.service.ProcessAnswers(
		c.Request.Context(),
		t,
		req.Answers,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
func (h *Handler) GetCandidates(c *gin.Context) {

	companyID := c.GetInt64("company_id")

	candidates, err := h.service.GetCandidates(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"candidates": candidates,
	})
}

func (h *Handler) GetCompanyMembers(c *gin.Context) {

	companyID := c.GetInt64("company_id")

	if companyID == 0 {
		c.JSON(400, gin.H{"error": "missing company"})
		return
	}

	members, err := h.service.GetCompanyMembers(
		c.Request.Context(),
		companyID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"members": members,
	})
}
func (h *Handler) GetCandidates2(c *gin.Context) {

	companyID := c.GetInt64("company_id")

	if companyID == 0 {
		c.JSON(400, gin.H{"error": "missing company"})
		return
	}

	candidates, err := h.service.GetCandidates2(
		c.Request.Context(),
		companyID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"candidates": candidates,
	})
}
