package main

import (
	"context"
	"log"
	"os"

	"github.com/alresave/mutari/internal/modules/matching"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

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

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := newDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	repo := matching.NewRepository(db)
	service := matching.NewService(repo)
	handler := matching.NewHandler(service)

	r := gin.Default()

	r.Static("/app", "./web")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/matches/:user_id", handler.GetMatches)

	log.Println("Mutari Engine Started")

	r.Run(":8080")
}
