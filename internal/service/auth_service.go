// internal/service/auth_service.go
package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/you/linkedinify/internal/config"
	"github.com/you/linkedinify/internal/model"
	"github.com/you/linkedinify/internal/repository"
)

type AuthService struct {
	repo repository.UserRepository
	key  []byte
}

func NewAuth(repo repository.UserRepository, cfg config.Config) *AuthService {
	return &AuthService{repo: repo, key: cfg.JWTSecret}
}

func (a *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		APIToken:     uuid.NewString(),
	}
	if err := a.repo.Create(ctx, user); err != nil {
		return "", err
	}
	return a.generateJWT(user.ID)
}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", jwt.ErrTokenInvalidAudience
	}
	return a.generateJWT(u.ID)
}

func (a *AuthService) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.key)
}
