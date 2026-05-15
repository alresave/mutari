package matching

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/alresave/mutari/internal/modules/models"
	"github.com/golang-jwt/jwt"
)

func vectorMagnitude(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v * v
	}
	return math.Sqrt(total)
}

func TestGetMCPQuestionsV1(t *testing.T) {
	questions := GetMCPQuestionsV1()

	if len(questions) != 64 {
		t.Fatalf("expected 64 questions, got %d", len(questions))
	}

	for i, q := range questions {
		if q.ID != i+1 {
			t.Fatalf("question %d has ID %d", i, q.ID)
		}
		if q.DimensionIndex != i {
			t.Fatalf("question %d has dimension index %d", i, q.DimensionIndex)
		}
		if q.Reverse {
			t.Fatalf("question %d should not be reverse", i)
		}
		if q.IsControl {
			t.Fatalf("question %d should not be control", i)
		}
	}
}

func TestBuildIdentityVectorNormalizesAndComputesConfidence(t *testing.T) {
	questions := []Question{
		{ID: 1, DimensionIndex: 0},
		{ID: 2, DimensionIndex: 1, Reverse: true},
		{ID: 3, DimensionIndex: 0, IsControl: true},
	}
	answers := []Answer{
		{QuestionID: 1, Value: 100},
		{QuestionID: 2, Value: 0},
		{QuestionID: 3, Value: 50},
	}

	vector, confidence := BuildIdentityVector(questions, answers)

	if len(vector) != 64 {
		t.Fatalf("expected 64 dimensions, got %d", len(vector))
	}

	if math.Abs(confidence-0.5) > 0.000001 {
		t.Fatalf("expected confidence 0.5, got %f", confidence)
	}

	magnitude := vectorMagnitude(vector)
	if math.Abs(magnitude-1) > 0.000001 {
		t.Fatalf("expected normalized vector, got magnitude %f", magnitude)
	}

	if vector[0] <= 0 {
		t.Fatalf("expected dimension 0 to remain positive, got %f", vector[0])
	}
	if vector[1] <= 0 {
		t.Fatalf("expected reversed dimension 1 to be positive, got %f", vector[1])
	}
}

func TestBuildVectorMapsAnswersByQuestionID(t *testing.T) {
	vector := BuildVector([]models.AnswerInput{
		{QuestionID: 3, Value: 5},
		{QuestionID: 1, Value: 2},
		{QuestionID: 65, Value: 4},
		{QuestionID: 2, Value: 1},
	})

	if len(vector) != 64 {
		t.Fatalf("expected 64 dimensions, got %d", len(vector))
	}

	if vector[0] != 2 || vector[1] != 1 || vector[2] != 5 {
		t.Fatalf("answers were not mapped by question ID: first values %#v", vector[:3])
	}
}

func TestPredictRole(t *testing.T) {
	tests := []struct {
		name   string
		blocks map[string]float64
		want   string
	}{
		{
			name: "execution leader",
			blocks: map[string]float64{
				"Ambition":            0.9,
				"Emotional Stability": 0.81,
			},
			want: "Execution Leader",
		},
		{
			name: "strategic architect",
			blocks: map[string]float64{
				"Cognitive Style": 0.89,
				"Growth":          0.81,
			},
			want: "Strategic Architect",
		},
		{
			name: "team integrator",
			blocks: map[string]float64{
				"Relational Architecture": 0.86,
				"Emotional Stability":     0.81,
			},
			want: "Team Integrator",
		},
		{
			name: "specialist",
			blocks: map[string]float64{
				"Cognitive Style": 0.91,
			},
			want: "Specialist",
		},
		{
			name: "balanced contributor",
			blocks: map[string]float64{
				"Ambition": 0.7,
			},
			want: "Balanced Contributor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PredictRole(tt.blocks); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestSuggestRoleUsesBlockAverages(t *testing.T) {
	tests := []struct {
		name  string
		block int
		want  string
	}{
		{name: "cognitive strategist", block: 1, want: "Strategist"},
		{name: "execution leader", block: 0, want: "Leader"},
		{name: "emotional connector", block: 3, want: "Connector"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vector := make([]float32, 64)
			for i := tt.block * 8; i < tt.block*8+8; i++ {
				vector[i] = 0.9
			}

			if got := suggestRole(vector); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGenerateTraits(t *testing.T) {
	traits := GenerateTraits(map[string]float64{
		"Cognitive Style":         0.9,
		"Emotional Stability":     0.81,
		"Ambition":                0.86,
		"Relational Architecture": 0.81,
		"Growth":                  0.81,
		"Core Values":             0.81,
	})

	if len(traits) != 6 {
		t.Fatalf("expected 6 traits, got %d", len(traits))
	}

	fallback := GenerateTraits(map[string]float64{
		"Ambition": 0.2,
	})
	if len(fallback) != 1 || fallback[0] != "Balanced team contributor" {
		t.Fatalf("unexpected fallback traits: %#v", fallback)
	}
}

func TestGenerateJWTWithCompany(t *testing.T) {
	tokenString, err := GenerateJWTWithCompany(42, 7)
	if err != nil {
		t.Fatalf("GenerateJWTWithCompany returned error: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		t.Fatalf("unable to parse token: %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("unexpected claims type: %T", token.Claims)
	}

	if claims["user_id"].(float64) != 42 {
		t.Fatalf("expected user_id 42, got %v", claims["user_id"])
	}
	if claims["company_id"].(float64) != 7 {
		t.Fatalf("expected company_id 7, got %v", claims["company_id"])
	}
	if claims["exp"].(float64) <= float64(jwt.TimeFunc().Unix()) {
		t.Fatalf("expected future exp, got %v", claims["exp"])
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if len(token) != 64 {
		t.Fatalf("expected 64-char token, got %d", len(token))
	}

	if _, err := hex.DecodeString(token); err != nil {
		t.Fatalf("expected hex token, got %q: %v", token, err)
	}
}

func TestRelationshipSimulationAndLabels(t *testing.T) {
	result := RunRelationshipSimulation()

	if result.ConflictScore != 0.8 || result.LifestyleScore != 0.9 || result.EmotionalScore != 0.75 || result.FutureScore != 0.85 || result.DecisionScore != 0.7 {
		t.Fatalf("unexpected simulation result: %#v", result)
	}

	tests := []struct {
		score float64
		want  string
	}{
		{0.9, "Exceptional"},
		{0.8, "Strong"},
		{0.6, "Moderate"},
		{0.59, "Uncertain"},
	}

	for _, tt := range tests {
		if got := StabilityLabel(tt.score); got != tt.want {
			t.Fatalf("for score %f expected %q, got %q", tt.score, tt.want, got)
		}
	}
}
