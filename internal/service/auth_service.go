package service

import (
	"ai-assistant/internal/config"
	"ai-assistant/internal/domain"
	"ai-assistant/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password, fullName string) (*domain.User, string, error)
	LoginUser(ctx context.Context, email, password string) (*domain.User, string, error)
	ParseToken(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error)
}
type authSvc struct {
	userRepo repository.UserRepository
	cfg      *config.Config // (!!!) ИСПРАВЛЕНИЕ: Было "Cofig"
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService { // (!!!) ИСПРАВЛЕНИЕ: Было "Cofig"
	return &authSvc{userRepo: userRepo, cfg: cfg}
}

func (s *authSvc) RegisterUser(ctx context.Context, email, password, fullName string) (*domain.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userRepo.CreateUser(ctx, email, string(hash), fullName)
	if err != nil {
		return nil, "", err
	}

	// TODO: Создать счет в account_repo

	token, err := s.generateToken(user.ID)
	return user, token, err
}

func (s *authSvc) LoginUser(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid password")
	}

	token, err := s.generateToken(user.ID)
	return user, token, err
}

func (s *authSvc) generateToken(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTLifetime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authSvc) ParseToken(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}