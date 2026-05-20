package matching

import "github.com/alresave/mutari/internal/modules/models"

// GetMCPQuestionsV1 returns the question definitions
// for MCP model version 1.
func GetMCPQuestionsV1() []Question {

	questions := make([]Question, 0, models.MCPVectorSize)

	// One primary question per dimension for MVP wiring
	for i := 0; i < models.MCPVectorSize; i++ {
		questions = append(questions, Question{
			ID:             i + 1,
			DimensionIndex: i,
			Reverse:        false,
			IsControl:      false,
		})
	}

	return questions
}
