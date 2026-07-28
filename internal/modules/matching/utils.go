package matching

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("mutari-secret")

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"company_id": int64(0),
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSecret)
}

// PredictRole infers the most suitable role for a candidate
func PredictRole(blocks map[string]float64) string {

	ambition := blocks["Ambition"]
	cognitive := blocks["Cognitive Style"]
	emotional := blocks["Emotional Stability"]
	relational := blocks["Relational Architecture"]
	growth := blocks["Growth"]

	if ambition > 0.85 && emotional > 0.80 {
		return "Execution Leader"
	}

	if cognitive > 0.88 && growth > 0.80 {
		return "Strategic Architect"
	}

	if relational > 0.85 && emotional > 0.80 {
		return "Team Integrator"
	}

	if cognitive > 0.90 {
		return "Specialist"
	}

	return "Balanced Contributor"
}

// GenerateTraits converts block scores into readable traits
func GenerateTraits(blocks map[string]float64) []string {

	traits := []string{}

	if blocks["Cognitive Style"] > 0.85 {
		traits = append(traits, "Strategic thinker")
	}

	if blocks["Emotional Stability"] > 0.80 {
		traits = append(traits, "High emotional regulation")
	}

	if blocks["Ambition"] > 0.85 {
		traits = append(traits, "Strong execution drive")
	}

	if blocks["Relational Architecture"] > 0.80 {
		traits = append(traits, "Collaborative mindset")
	}

	if blocks["Growth"] > 0.80 {
		traits = append(traits, "Continuous improvement oriented")
	}

	if blocks["Core Values"] > 0.80 {
		traits = append(traits, "Values aligned with team culture")
	}

	if len(traits) == 0 {
		traits = append(traits, "Balanced team contributor")
	}

	return traits
}
func predictRole(blocks []float64) string {

	if blocks[1] > 0.85 && blocks[4] > 0.8 {
		return "Visionary"
	}

	if blocks[0] > 0.8 && blocks[6] > 0.75 {
		return "Operator"
	}

	if blocks[5] > 0.85 && blocks[7] > 0.8 {
		return "Leader"
	}

	return "Specialist"
}
