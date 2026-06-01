package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret         []byte
	accessTokenTTL time.Duration
	issuer         string
}

func NewJWTService(secret string, accessTokenTTL time.Duration, issuer string) *JWTService {
	return &JWTService{
		secret:         []byte(secret),
		accessTokenTTL: accessTokenTTL,
		issuer:         issuer,
	}
}

func (s *JWTService) GenerateAccessToken(userID string) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return signedToken, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secret, nil
		},
		jwt.WithIssuer(s.issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt")
	}

	return claims, nil
}
