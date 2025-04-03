package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var Conn *pgxpool.Pool

func ConnectDB() (*pgxpool.Pool, error) {
	const dbURL = "postgres://postgres:root@localhost:5432/db100"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("faied connect : %w", err)
	}
	log.Println("✅ Connected to PostgreSQL successfully!")
	Conn = pool
	return pool, nil
}
