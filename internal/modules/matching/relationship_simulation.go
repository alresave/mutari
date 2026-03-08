package matching

type SimulationResult struct {
	ConflictScore  float64
	LifestyleScore float64
	EmotionalScore float64
	FutureScore    float64
	DecisionScore  float64
}

func RunRelationshipSimulation() SimulationResult {

	return SimulationResult{
		ConflictScore:  0.8,
		LifestyleScore: 0.9,
		EmotionalScore: 0.75,
		FutureScore:    0.85,
		DecisionScore:  0.7,
	}
}
func StabilityLabel(score float64) string {

	switch {
	case score >= 0.9:
		return "Exceptional"

	case score >= 0.8:
		return "Strong"

	case score >= 0.6:
		return "Moderate"

	default:
		return "Uncertain"
	}
}
