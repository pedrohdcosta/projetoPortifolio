package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Service struct {
	repo      Repository
	jwtSecret []byte
}

type Repository interface {
	CreateUser(ctx context.Context, name, email, passHash string) (int64, error)
	FindUserByEmail(ctx context.Context, email string) (id int64, name, emailDB, passHash string, err error)
	FindUserByID(ctx context.Context, id int64) (User, error)
}

func NewService(r Repository, jwtSecret []byte) *Service {
	return &Service{repo: r, jwtSecret: jwtSecret}
}

func (s *Service) Signup(ctx context.Context, name, email, password string) (User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	id, err := s.repo.CreateUser(ctx, name, email, string(hash))
	if err != nil {
		return User{}, err
	}
	return User{ID: id, Name: name, Email: email}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, User, error) {
	id, name, emailDB, passHash, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", User{}, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password)) != nil {
		return "", User{}, errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{"sub": id, "email": emailDB, "exp": time.Now().Add(24 * time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(s.jwtSecret)
	return signed, User{ID: id, Name: name, Email: emailDB}, nil
}
