package matching

import (
	"math"

	"github.com/alresave/mutari/internal/modules/models"
)

// Question describes how a quiz question maps into the 64-dimensional MCP vector.
type Question struct {
	ID             int
	DimensionIndex int
	Reverse        bool
	IsControl      bool
}

// Answer stores the raw score submitted for a quiz question.
type Answer struct {
	QuestionID int
	Value      float64
}

const dimensionSize = models.MCPVectorSize

// contrastGamma controls how strongly scores move away from the neutral midpoint.
const contrastGamma = 1.8

// orthogonalMix reduces shared signal across dimensions before final normalization.
const orthogonalMix = 0.15

// BuildIdentityVector converts question answers into a normalized identity vector.
// It also returns a confidence score based on the consistency of control questions.
func BuildIdentityVector(
	questions []Question,
	answers []Answer,
) ([]float64, float64) {

	vector := make([]float64, dimensionSize)
	counts := make([]int, dimensionSize)

	answerMap := make(map[int]float64)

	for _, a := range answers {

		v := a.Value / 100.0

		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}

		answerMap[a.QuestionID] = v
	}

	var inconsistency float64
	var controlChecks int

	for _, q := range questions {

		val, ok := answerMap[q.ID]
		if !ok {
			continue
		}

		if q.Reverse {
			val = 1 - val
		}

		dim := q.DimensionIndex

		if dim < 0 || dim >= dimensionSize {
			continue
		}

		if q.IsControl {

			if counts[dim] > 0 {

				current := vector[dim] / float64(counts[dim])
				diff := math.Abs(current - val)

				inconsistency += diff
				controlChecks++
			}

			continue
		}

		vector[dim] += val
		counts[dim]++
	}

	// Average each populated dimension, then center scores around zero so cosine
	// similarity compares direction rather than raw positive-only magnitudes.
	for i := 0; i < dimensionSize; i++ {

		if counts[i] > 0 {
			vector[i] /= float64(counts[i])
		}

		centered := vector[i] - 0.5

		if centered >= 0 {
			centered = math.Pow(centered, contrastGamma)
		} else {
			centered = -math.Pow(math.Abs(centered), contrastGamma)
		}

		vector[i] = centered
	}

	// Subtract a small amount of the other dimensions' mean to reduce broad,
	// non-specific agreement from dominating the final profile.
	for i := 0; i < dimensionSize; i++ {

		var mix float64

		for j := 0; j < dimensionSize; j++ {

			if i == j {
				continue
			}

			mix += vector[j]
		}

		mix /= float64(dimensionSize - 1)

		vector[i] = vector[i] - orthogonalMix*mix
	}

	// Normalize to unit length so downstream cosine comparisons are stable.
	var mag float64

	for _, v := range vector {
		mag += v * v
	}

	mag = math.Sqrt(mag)

	if mag > 0 {

		for i := range vector {
			vector[i] /= mag
		}
	}

	// Confidence starts perfect and is reduced by disagreement with control
	// questions that target dimensions already answered by primary questions.
	confidence := 1.0

	if controlChecks > 0 {

		avgInconsistency := inconsistency / float64(controlChecks)

		confidence = 1 - avgInconsistency

		if confidence < 0 {
			confidence = 0
		}
	}

	return vector, confidence
}
