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

		tokenString := strings.Replace(header, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		log.Println("**************************jwt debug******************************")

		claims := token.Claims.(jwt.MapClaims)
		fmt.Printf("Claims: %v\n", claims)

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
	r.Static("/app", "./web")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	r.POST("/auth/signup", handler.Signup)
	r.POST("/auth/login", handler.Login)

	// Public endpoint
	r.GET("/mcp/questions", handler.GetMCPQuestions)

	// Protected routes
	auth := r.Group("/")
	auth.Use(AuthMiddleware())

	auth.GET("/matches", handler.GetMatches)
	r.POST("/admin/rebuild-profiles", handler.RebuildProfiles)
	auth.POST("/mcp/answers", handler.SubmitMCPAnswers)

	log.Println("Mutari Engine Started")
	log.Println("**************************jwt debug******************************")

	r.Run(":8080")
}
