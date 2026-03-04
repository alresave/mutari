package matching

import (
	"math"
)

//
// ----------------------
// DOMAIN MODELS
// ----------------------
//

// Question defines how a test question maps into the identity vector.
type Question struct {
	ID             int  // Unique question identifier
	DimensionIndex int  // 0–63
	Reverse        bool // If true, invert value (control polarity)
	IsControl      bool // If true, used only for consistency scoring
}

// Answer represents a user response.
type Answer struct {
	QuestionID int
	Value      float64 // Expected 0–100 slider input
}

//
// ----------------------
// IDENTITY VECTOR BUILDER
// ----------------------
//

// BuildIdentityVector builds the 64-d identity vector
// and calculates a structural confidence score (0–1).
func BuildIdentityVector(
	questions []Question,
	answers []Answer,
) ([]float64, float64) {

	const dimensionSize = 64

	vector := make([]float64, dimensionSize)
	counts := make([]int, dimensionSize)

	// Map answers by question ID
	answerMap := make(map[int]float64)
	for _, a := range answers {
		normalized := a.Value / 100.0
		if normalized < 0 {
			normalized = 0
		}
		if normalized > 1 {
			normalized = 1
		}
		answerMap[a.QuestionID] = normalized
	}

	var inconsistency float64
	var controlChecks int

	for _, q := range questions {

		val, ok := answerMap[q.ID]
		if !ok {
			continue
		}

		// Apply reverse polarity if needed
		if q.Reverse {
			val = 1.0 - val
		}

		dim := q.DimensionIndex

		if dim < 0 || dim >= dimensionSize {
			continue // Safety guard
		}

		if q.IsControl {
			// Compare against accumulated primary dimension value
			if counts[dim] > 0 {
				currentAvg := vector[dim] / float64(counts[dim])
				diff := math.Abs(currentAvg - val)
				inconsistency += diff
				controlChecks++
			}
			continue
		}

		// Primary structural contribution
		vector[dim] += val
		counts[dim]++
	}

	// Average each dimension
	for i := 0; i < dimensionSize; i++ {
		if counts[i] > 0 {
			vector[i] /= float64(counts[i])
		}
	}

	// Confidence calculation
	confidence := 1.0
	if controlChecks > 0 {
		avgInconsistency := inconsistency / float64(controlChecks)
		confidence = 1.0 - avgInconsistency
		if confidence < 0 {
			confidence = 0
		}
	}

	return vector, confidence
}
