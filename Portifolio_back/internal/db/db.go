package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	return pgxpool.New(ctx, url)
}
