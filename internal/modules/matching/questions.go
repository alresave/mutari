package matching

// GetMCPQuestionsV1 returns the question definitions
// for MCP model version 1.
func GetMCPQuestionsV1() []Question {

	const dimensionSize = 64

	questions := make([]Question, 0, dimensionSize)

	// One primary question per dimension for MVP wiring
	for i := 0; i < dimensionSize; i++ {
		questions = append(questions, Question{
			ID:             i + 1,
			DimensionIndex: i,
			Reverse:        false,
			IsControl:      false,
		})
	}

	return questions
}
