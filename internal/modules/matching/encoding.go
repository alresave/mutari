package matching

import (
	"math"
)

type Question struct {
	ID             int
	DimensionIndex int
	Reverse        bool
	IsControl      bool
}

type Answer struct {
	QuestionID int
	Value      float64
}

const dimensionSize = 64

// dispersion exponent (controla qué tan extremos se vuelven los perfiles)
const contrastGamma = 1.8

// mixing factor para ortogonalizar dimensiones
const orthogonalMix = 0.15

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

	// average dimensions
	for i := 0; i < dimensionSize; i++ {

		if counts[i] > 0 {
			vector[i] /= float64(counts[i])
		}

		// center distribution around 0
		centered := vector[i] - 0.5

		// amplify contrast
		if centered >= 0 {
			centered = math.Pow(centered, contrastGamma)
		} else {
			centered = -math.Pow(math.Abs(centered), contrastGamma)
		}

		vector[i] = centered
	}

	// ----------------------
	// ORTHOGONAL MIXING
	// ----------------------

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

	// ----------------------
	// L2 NORMALIZATION
	// ----------------------

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

	// ----------------------
	// CONFIDENCE SCORE
	// ----------------------

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
