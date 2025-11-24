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
	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/devices"
	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/telemetry"
)

func main() {
	_ = godotenv.Load("../../.env.example")
	ctx := context.Background()
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}

	r := setupRouter(ctx, pool)
	port := getPort()

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// setupRouter configures and returns the Gin router with all routes and middleware.
func setupRouter(ctx context.Context, pool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint
	r.GET("/health", healthCheckHandler)

	// Database schema initialization
	ensureSchema(ctx, pool)

	// Wrap pool for interfaces
	wrapped := wrap(pool)

	// API routes
	auth.RegisterRoutes(r, wrapped)

	// Protected API routes (require authentication)
	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		// Devices CRUD
		devicesRepo := devices.NewRepo(&devicesQuerier{wrapped})
		devicesHandler := devices.NewHandler(devicesRepo)
		devicesHandler.RegisterRoutes(api)

		// Telemetry CRUD
		telemetryRepo := telemetry.NewRepo(&telemetryQuerier{wrapped})
		telemetryHandler := telemetry.NewHandler(telemetryRepo)
		telemetryHandler.RegisterRoutes(api)
	}

	// Configure static file serving for frontend SPA
	configureStaticFiles(r)

	return r
}

// healthCheckHandler returns a simple health check response.
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// configureStaticFiles sets up static file serving and SPA routing.
// In production/Docker deployments, the frontend is built and served from ./static
func configureStaticFiles(r *gin.Engine) {
	const staticDir = "./static"

	stat, err := os.Stat(staticDir)
	if err != nil {
		log.Printf("Static directory not found (%v), serving API only", err)
		r.NoRoute(apiOnlyNoRouteHandler)
		return
	}

	if !stat.IsDir() {
		log.Printf("Static path exists but is not a directory, serving API only")
		r.NoRoute(apiOnlyNoRouteHandler)
		return
	}

	log.Printf("Serving static files from %s", staticDir)
	r.Static("/assets", staticDir+"/assets")
	r.StaticFile("/", staticDir+"/index.html")
	r.StaticFile("/favicon.ico", staticDir+"/favicon.ico")
	r.NoRoute(spaNoRouteHandler(staticDir))
}

// spaNoRouteHandler returns a handler for SPA client-side routing.
// API routes get JSON 404, all other routes serve index.html.
func spaNoRouteHandler(staticDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// API routes should return JSON 404, not the SPA
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			log.Printf("API route not found: %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		// Serve index.html for all other routes (SPA client-side routing)
		c.File(staticDir + "/index.html")
	}
}

// apiOnlyNoRouteHandler returns JSON 404 for all routes when no static files are served.
func apiOnlyNoRouteHandler(c *gin.Context) {
	log.Printf("Route not found: %s %s", c.Request.Method, c.Request.URL.Path)
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

// getPort determines the port to listen on from environment variables.
// Priority: PORT (Azure) > APP_PORT > default "8080"
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		return port
	}
	return "8080"
}

type pgxWrap struct{ *pgxpool.Pool }

func wrap(p *pgxpool.Pool) *pgxWrap { return &pgxWrap{p} }
func (w *pgxWrap) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := w.Pool.Exec(ctx, sql, args...)
	return err
}
func (w *pgxWrap) QueryRow(ctx context.Context, sql string, args ...any) interface{ Scan(dest ...any) error } {
	return w.Pool.QueryRow(ctx, sql, args...)
}

// pgxRows wraps pgx.Rows to satisfy various Rows interfaces.
type pgxRows struct {
	rows interface {
		Next() bool
		Scan(dest ...any) error
		Close()
		Err() error
	}
}

func (r *pgxRows) Next() bool            { return r.rows.Next() }
func (r *pgxRows) Scan(dest ...any) error { return r.rows.Scan(dest...) }
func (r *pgxRows) Close()                { r.rows.Close() }
func (r *pgxRows) Err() error            { return r.rows.Err() }

// Query returns a Rows-compatible result for devices package.
func (w *pgxWrap) Query(ctx context.Context, sql string, args ...any) (devices.Rows, error) {
	r, err := w.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: r}, nil
}

// devicesQuerier adapts pgxWrap to devices.RowsQuerier interface.
type devicesQuerier struct{ *pgxWrap }

func (d *devicesQuerier) Query(ctx context.Context, sql string, args ...any) (devices.Rows, error) {
	return d.pgxWrap.Query(ctx, sql, args...)
}

// telemetryQuerier adapts pgxWrap to telemetry.RowsQuerier interface.
type telemetryQuerier struct{ *pgxWrap }

func (t *telemetryQuerier) Query(ctx context.Context, sql string, args ...any) (telemetry.Rows, error) {
	r, err := t.pgxWrap.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: r}, nil
}

func ensureSchema(ctx context.Context, p *pgxpool.Pool) {
	// Create users table
	_, _ = p.Exec(ctx, `CREATE TABLE IF NOT EXISTS app_user(
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)`)

	// Create devices table
	_, _ = p.Exec(ctx, `CREATE TABLE IF NOT EXISTS device(
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		room TEXT,
		type TEXT DEFAULT 'smart_plug',
		status TEXT DEFAULT 'offline',
		metadata TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		last_seen TIMESTAMPTZ
	)`)

	// Create index on user_id for faster device lookups
	_, _ = p.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_device_user_id ON device(user_id)`)

	// Create telemetry table
	_, _ = p.Exec(ctx, `CREATE TABLE IF NOT EXISTS telemetry(
		id BIGSERIAL PRIMARY KEY,
		device_id BIGINT NOT NULL REFERENCES device(id) ON DELETE CASCADE,
		power DOUBLE PRECISION NOT NULL,
		voltage DOUBLE PRECISION,
		current DOUBLE PRECISION,
		timestamp TIMESTAMPTZ DEFAULT NOW()
	)`)

	// Create index on device_id and timestamp for faster telemetry queries
	_, _ = p.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_telemetry_device_id ON telemetry(device_id)`)
	_, _ = p.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_telemetry_timestamp ON telemetry(timestamp DESC)`)
}
