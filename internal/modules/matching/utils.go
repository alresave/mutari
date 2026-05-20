package matching

import (
	"crypto/rand"
	"os"
	"time"

	"encoding/hex"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/golang-jwt/jwt"
)

const defaultJWTSecret = "mutari-secret"

func jwtSigningSecret() []byte {
	if secret := os.Getenv("MUTARI_JWT_SECRET"); secret != "" {
		return []byte(secret)
	}

	return []byte(defaultJWTSecret)
}

// GenerateJWT creates a JWT token for a user.
func GenerateJWT(userID int64) (string, error) {
	return GenerateJWTWithCompany(userID, 0)
}

// GenerateJWTWithCompany creates a JWT token that includes the current company context.
func GenerateJWTWithCompany(userID, companyID int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"company_id": companyID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSigningSecret())
}

// PredictRole infers the most suitable role for a candidate
func PredictRole(blocks map[string]float64) string {

	ambition := scoreForBlock(blocks, "Learning Velocity", "Ambition")
	cognitive := scoreForBlock(blocks, "Cognitive Processing", "Cognitive Style")
	emotional := scoreForBlock(blocks, "Emotional Resilience", "Emotional Stability")
	relational := scoreForBlock(blocks, "Collaboration Depth", "Relational Architecture")
	growth := scoreForBlock(blocks, "Team Positioning", "Growth")

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

	if scoreForBlock(blocks, "Cognitive Processing", "Cognitive Style") > 0.85 {
		traits = append(traits, "Strategic thinker")
	}

	if scoreForBlock(blocks, "Emotional Resilience", "Emotional Stability") > 0.80 {
		traits = append(traits, "High emotional regulation")
	}

	if scoreForBlock(blocks, "Learning Velocity", "Ambition") > 0.85 {
		traits = append(traits, "Strong execution drive")
	}

	if scoreForBlock(blocks, "Collaboration Depth", "Relational Architecture") > 0.80 {
		traits = append(traits, "Collaborative mindset")
	}

	if scoreForBlock(blocks, "Team Positioning", "Growth") > 0.80 {
		traits = append(traits, "Continuous improvement oriented")
	}

	if scoreForBlock(blocks, "Professional Values", "Core Values") > 0.80 {
		traits = append(traits, "Values aligned with team culture")
	}

	if len(traits) == 0 {
		traits = append(traits, "Balanced team contributor")
	}

	return traits
}

func scoreForBlock(blocks map[string]float64, names ...string) float64 {
	for _, name := range names {
		if score, ok := blocks[name]; ok {
			return score
		}
	}

	return 0
}

func predictRole(blocks []float64) string {
	if len(blocks) < models.MCPBlockCount {
		return "Specialist"
	}

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

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
