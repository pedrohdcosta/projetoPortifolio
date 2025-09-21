package auth

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type repoPG struct{ q querier }

type querier interface {
	Exec(ctx context.Context, sql string, args ...any) error
	QueryRow(ctx context.Context, sql string, args ...any) row
}

type row interface{ Scan(dest ...any) error }

func (r *repoPG) CreateUser(ctx context.Context, name, email, passHash string) (int64, error) {
	sql := `insert into app_user(name,email,password_hash) values($1,$2,$3) returning id`
	var id int64
	err := r.q.QueryRow(ctx, sql, name, email, passHash).Scan(&id)
	return id, err
}
func (r *repoPG) FindUserByEmail(ctx context.Context, email string) (int64, string, string, string, error) {
	sql := `select id,name,email,password_hash from app_user where email=$1`
	var id int64
	var name, emailDB, pass string
	err := r.q.QueryRow(ctx, sql, email).Scan(&id, &name, &emailDB, &pass)
	return id, name, emailDB, pass, err
}
func (r *repoPG) FindUserByID(ctx context.Context, id int64) (User, error) {
	sql := `select id,name,email from app_user where id=$1`
	var u User
	err := r.q.QueryRow(ctx, sql, id).Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

func RegisterRoutes(r *gin.Engine, q querier) {
	repo := &repoPG{q: q}
	svc := NewService(repo, []byte(os.Getenv("JWT_SECRET")))
	g := r.Group("/api/auth")
	g.POST("/signup", func(c *gin.Context) {
		var in struct{ Name, Email, Password string }
		if err := c.BindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		u, err := svc.Signup(c, in.Name, in.Email, in.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, u)
	})
	g.POST("/login", func(c *gin.Context) {
		var in struct{ Email, Password string }
		if err := c.BindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		t, u, err := svc.Login(c, in.Email, in.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"accessToken": t, "user": u})
	})
	g.GET("/me", AuthMiddleware(), func(c *gin.Context) {
		idStr := c.GetString("sub")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		u, err := repo.FindUserByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	})
}
