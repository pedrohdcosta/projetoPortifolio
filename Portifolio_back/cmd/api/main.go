package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/auth"
	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/db"
)

func main() {
	ctx := context.Background()
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	// migração mínima
	ensureSchema(ctx, pool)
	// rotas de auth
	auth.RegisterRoutes(r, wrap(pool))
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

type pgxWrap struct{ *pgxpool.Pool }

func wrap(p *pgxpool.Pool) pgxWrap { return pgxWrap{p} }
func (w pgxWrap) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := w.Pool.Exec(ctx, sql, args...)
	return err
}
func (w pgxWrap) QueryRow(ctx context.Context, sql string, args ...any) dbRow {
	return w.Pool.QueryRow(ctx, sql, args...)
}

type dbRow interface{ Scan(dest ...any) error }

func ensureSchema(ctx context.Context, p *pgxpool.Pool) {
	_, _ = p.Exec(ctx, `create table if not exists app_user(
	id bigserial primary key,
	name text not null,
	email text unique not null,
	password_hash text not null,
	created_at timestamptz default now()
	)`)
}
