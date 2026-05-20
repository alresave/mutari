package mcp

import (
	"context"

	"github.com/alresave/mutari/internal/modules/models"
)

// RepositoryStore defines MCP question, answer, and profile persistence behavior.
type RepositoryStore interface {
	GetMCPQuestions(ctx context.Context) ([]models.MCPQuestion, error)
	StoreCompanyMemberAnswers(ctx context.Context, memberID int64, answers []models.MCPAnswer) error
	BuildCompanyMemberProfile(ctx context.Context, memberID int64) error
	HasMCPProfile(ctx context.Context, userID int64) (bool, error)
	GetAllCompanyMemberIDs(ctx context.Context) ([]int64, error)
	SaveAnswers(ctx context.Context, candidateID int64, answers []models.Answer) error
	SaveMemberAnswers(ctx context.Context, memberID int64, answers []models.AnswerInput) error
	UpsertCandidateProfile(ctx context.Context, candidateID int64, vector []float32) error
	UpsertCompanyMemberProfile(ctx context.Context, memberID int64, vector []float32) error
}

var _ RepositoryStore = (*Repository)(nil)
