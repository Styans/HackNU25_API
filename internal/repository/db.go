package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, databaseURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// if err = pool.Ping(ctx); err != nil {
	// 	log.Fatalf("Failed to ping database: %v\n", err)
	// }

	log.Println("Database connection established successfully")
	return pool
}