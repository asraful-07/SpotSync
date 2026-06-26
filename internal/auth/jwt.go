package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//  Removed hardcoded jwtSecretKey and defaultTokenDuration constants
// — config should be the single source of truth

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userID uint, name, email, role, phone, createdAt string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string, tokenDuration time.Duration) JWTService {
	return &jwtService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (s *jwtService) GenerateToken(userID uint, name, email, role, phone, createdAt string) (string, error) {
	claims := &JWTClaims{
		UserID:    userID,
		Name:      name,
		Email:     email,
		Role:      role,
		Phone:     phone,
		CreatedAt: createdAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}