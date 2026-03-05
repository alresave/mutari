package matching

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}
type MCPAnswer struct {
	QuestionID int     `json:"question_id"`
	Block      int     `json:"block"`
	Value      float64 `json:"value"`
}

type SubmitAnswersRequest struct {
	UserID  int64       `json:"user_id"`
	Answers []MCPAnswer `json:"answers"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMatches(c *gin.Context) {

	userIDParam := c.Param("user_id")

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	matches, err := h.service.GetMatches(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"matches": matches,
	})
}
func (h *Handler) SubmitMCPAnswers(c *gin.Context) {

	var req SubmitAnswersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.SubmitMCPAnswers(c.Request.Context(), req.UserID, req.Answers)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "answers stored",
	})
}

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
