package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/permissions"
	"github.com/you/pawtrack/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface for authentication business logic
type AuthService interface {
	Login(email, password string) (string, *models.User, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// TokenClaims custom JWT claims
type TokenClaims struct {
	UserID uint            `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// authService implementation of the auth service
type authService struct {
	userRepo  repository.UserRepository
	permRepo  repository.PermissionRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, permRepo repository.PermissionRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-me" // Fallback for dev
	}

	expiryHours := 24
	if envExpiry := os.Getenv("JWT_EXPIRY_HOURS"); envExpiry != "" {
		// Parse or use default
		expiryHours = 24
	}

	return &authService{
		userRepo:  userRepo,
		permRepo:  permRepo,
		jwtSecret: []byte(secret),
		jwtExpiry: time.Duration(expiryHours) * time.Hour,
	}
}

// Login authenticates user and returns JWT token
func (s *authService) Login(email, password string) (string, *models.User, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// ValidateToken validates JWT token and returns claims
func (s *authService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// generateToken creates a new JWT token for user
func (s *authService) generateToken(user *models.User) (string, error) {
	claims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
