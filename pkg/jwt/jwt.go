package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	errSignJwtToken    = errors.New("failed to sign JWT token")
	errInvalidJwtToken = errors.New("invalid JWT token")
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type Jwt struct {
	key        string
	expiration time.Duration
}

func New(key string, expiration time.Duration) *Jwt {
	return &Jwt{key: key, expiration: expiration}
}

func (j *Jwt) GenerateToken(userID, email, role string) (string, error) {
	now := time.Now()
	claims := &CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(j.expiration).Unix(),
			IssuedAt:  now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.key))
	if err != nil {
		return "", errSignJwtToken
	}
	return tokenString, nil
}

func (j *Jwt) ValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidJwtToken
		}
		return []byte(j.key), nil
	})
	if err != nil {
		return nil, errInvalidJwtToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errInvalidJwtToken
	}
	return claims, nil
}
