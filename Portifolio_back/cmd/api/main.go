package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/auth"
	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/db"
)

func main() {
	_ = godotenv.Load(".env")
	ctx := context.Background()
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	r.Use(gin.Logger(), gin.Recovery())

	// health rápido
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// loga o que vier errado
	r.NoRoute(func(c *gin.Context) {
		log.Printf("NoRoute %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

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
func (w pgxWrap) QueryRow(ctx context.Context, sql string, args ...any) interface{ Scan(dest ...any) error } {
	return w.Pool.QueryRow(ctx, sql, args...)
}

func ensureSchema(ctx context.Context, p *pgxpool.Pool) {
	_, _ = p.Exec(ctx, `create table if not exists app_user(
	id bigserial primary key,
	name text not null,
	email text unique not null,
	password_hash text not null,
	created_at timestamptz default now()
	)`)
}
