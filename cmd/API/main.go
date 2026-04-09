package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alresave/mutari/internal/modules/matching"
	"github.com/alresave/mutari/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

var jwtSecret = []byte("mutari-secret")

func newDB(connString string) (*pgxpool.Pool, error) {

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	return pgxpool.NewWithConfig(ctx, config)
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.Replace(header, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			// 🔒 validar algoritmo
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
			return
		}

		// 🔒 expiration
		if expVal, ok := claims["exp"]; ok {
			if expFloat, ok := expVal.(float64); ok {
				if time.Now().Unix() > int64(expFloat) {
					c.AbortWithStatusJSON(401, gin.H{"error": "token expired"})
					return
				}
			}
		}

		// ✅ user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid user_id"})
			return
		}
		userID := int64(userIDFloat)

		// ✅ company_id seguro
		var companyID int64
		if val, ok := claims["company_id"]; ok {
			if f, ok := val.(float64); ok {
				companyID = int64(f)
			}
		}

		c.Set("user_id", userID)
		c.Set("company_id", companyID)

		c.Next()
	}
}

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := newDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	securityLimiter := security.NewRateLimiter()

	repo := matching.NewRepository(db)
	service := matching.NewService(repo)
	handler := matching.NewHandler(service, securityLimiter)

	r := gin.Default()

	// Static UI
	r.Static("app", "web/app")
	r.Static("dating", "web/dating")

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth endpoints
	r.POST("/auth/signup", handler.Signup)
	r.POST("/auth/login", handler.Login)

	// Public MCP endpoint
	r.GET("/mcp/questions", handler.GetMCPQuestions)

	// Protected routes
	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	auth.GET("/me", handler.Me)

	auth.POST("/mcp/answers", handler.SubmitAnswers)

	auth.GET("/matches", handler.GetMatches)
	auth.GET("/teams/:team_id/fit/:user_id", handler.GetTeamMemberFit)
	auth.GET("/teams/:team_id/fit/member/:member_id", handler.GetTeamMemberFit)
	auth.GET("/teams/:team_id/simulate/:candidate_id", handler.SimulateTeamFit)
	auth.GET("/teams/:team_id/culture", handler.GetTeamCulture)
	auth.GET("/teams/:team_id/gaps", handler.GetTeamGaps)
	auth.GET("/teams/:team_id/vector", handler.GetTeamVector)
	auth.GET("/teams", handler.GetTeams)
	auth.GET("/teams/:team_id/hiring-profile", handler.GetTeamHiringProfile)
	auth.POST("/companies", handler.CreateCompany)
	auth.GET("/company-members", handler.GetCompanyMembers)
	auth.POST("/company-members", handler.CreateCompanyMember)
	auth.POST("/tokens", handler.GenerateToken)
	auth.POST("/candidates", handler.CreateCandidate)
	auth.GET("/candidates", handler.GetCandidates2)
	r.GET("/quiz", handler.GetQuiz)

	auth.POST("/teams/:team_id/members", handler.AddMemberToTeam)
	auth.POST("/teams", handler.CreateTeam)

	// TEAM ENGINE
	auth.GET("/teams/:team_id/matches", handler.GetTeamMatches)

	// Admin utility
	r.POST("/admin/rebuild-profiles", handler.RebuildProfiles)
	auth.POST("/admin/rebuild-team-profiles", handler.RebuildTeamProfiles)

	log.Println("Mutari Engine Started")
	for _, r := range r.Routes() {
		fmt.Println(r.Method, r.Path)
	}

	r.Run(":8080")
}
