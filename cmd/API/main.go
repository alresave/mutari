package main

import (
	"context"
	"log"
	"os"
	"strings"

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
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := int64(claims["user_id"].(float64))

		c.Set("userID", userID)

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

	auth.POST("/mcp/answers", handler.SubmitMCPAnswers)

	auth.GET("/matches", handler.GetMatches)
	auth.GET("/teams/:team_id/fit/:user_id", handler.GetTeamFit)
	auth.GET("/teams/:team_id/simulate/:candidate_id", handler.SimulateTeamFit)
	auth.GET("/teams/:team_id/culture", handler.GetTeamCulture)
	auth.GET("/teams/:team_id/gaps", handler.GetTeamGaps)
	auth.GET("/teams", handler.GetTeams)

	// TEAM ENGINE
	auth.GET("/teams/:team_id/matches", handler.GetTeamMatches)

	// Admin utility
	r.POST("/admin/rebuild-profiles", handler.RebuildProfiles)

	log.Println("Mutari Engine Started")

	r.Run(":8080")
}
