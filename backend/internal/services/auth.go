package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type AuthService struct {
	db        *bun.DB
	jwtSecret []byte
	expiry    time.Duration
}

func NewAuthService(db *bun.DB, secret string, expiryHours int) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(secret),
		expiry:    time.Duration(expiryHours) * time.Hour,
	}
}

type LoginResult struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	var user models.User
	err := s.db.NewSelect().Model(&user).
		Where(`"email" = ?`, email).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("authService.Login: user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("authService.Login: invalid credentials")
	}
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("authService.Login: %w", err)
	}
	return &LoginResult{Token: token, User: &user}, nil
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("authService.Register: name is required")
	}
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("authService.Register: email is required")
	}
	if !isValidEmail(email) {
		return nil, fmt.Errorf("authService.Register: invalid email format")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("authService.Register: password must be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("authService.Register: hash: %w", err)
	}
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}
	if _, err := s.db.NewInsert().Model(user).Exec(ctx); err != nil {
		return nil, fmt.Errorf("authService.Register: insert: %w", err)
	}
	return user, nil
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) generateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	atIndex := strings.Index(email, "@")
	if atIndex <= 0 || atIndex == len(email)-1 {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if strings.Contains(parts[0], " ") || strings.Contains(parts[1], " ") {
		return false
	}
	return true
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("authService.ValidateToken: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("authService.ValidateToken: invalid token")
	}
	return claims, nil
}
