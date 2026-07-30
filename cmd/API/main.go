package main

import (
	"context"
	"fmt"
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

		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
			return
		}

		userID, ok := numberClaim(claims, "user_id")
		if !ok || userID <= 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
			return
		}
		companyID, _ := numberClaim(claims, "company_id")

		c.Set("userID", userID)
		c.Set("companyID", companyID)

		c.Next()
	}
}

func numberClaim(claims jwt.MapClaims, name string) (int64, bool) {
	value, exists := claims[name]
	if !exists {
		return 0, false
	}

	number, ok := value.(float64)
	return int64(number), ok && number == float64(int64(number))
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
	r.StaticFile("/swagger/doc.json", "docs/swagger.json")
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>Mutari API documentation</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head>
<body><div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>SwaggerUIBundle({url: '/swagger/doc.json', dom_id: '#swagger-ui'})</script>
</body></html>`))
	})

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
	auth.GET("/teams/:team_id/vector", handler.GetTeamVector)
	auth.GET("/teams", handler.GetTeams)
	auth.GET("/teams/:team_id/hiring-profile", handler.GetTeamHiringProfile)

	// TEAM ENGINE
	auth.GET("/teams/:team_id/matches", handler.GetTeamMatches)

	// Administrative utilities require authentication.
	auth.POST("/admin/rebuild-profiles", handler.RebuildProfiles)
	auth.POST("/admin/rebuild-team-profiles", handler.RebuildTeamProfiles)

	log.Println("Mutari Engine Started")

	r.Run(":8080")
}
