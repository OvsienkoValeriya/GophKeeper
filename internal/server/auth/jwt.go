package auth

import (
	"fmt"
	"time"

	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

func NewJWTConfig(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTConfig {
	return &JWTConfig{
		SecretKey:       secretKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

// GenerateJWT generates a JWT for the user
// Parameters:
//   - user: user
//
// Returns:
//   - string: JWT
//   - error: error if the JWT generation failed
func (config *JWTConfig) GenerateJWT(user *models.User) (string, error) {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Username: user.Username,
		UserID:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// GenerateRefreshToken generates a refresh token for the user
// Parameters:
//   - user: user
//
// Returns:
//   - string: refresh token
//   - error: error if the refresh token generation failed
func (config *JWTConfig) GenerateRefreshToken(user *models.User) (string, error) {
	claims := &RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey + "-refresh"))
}

// VerifyRefreshToken verifies a refresh token
// Parameters:
//   - tokenString: refresh token
//
// Returns:
//   - *RefreshClaims: refresh claims
//   - error: error if the refresh token verification failed
func (config *JWTConfig) VerifyRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey + "-refresh"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

// VerifyToken verifies a JWT
// Parameters:
//   - tokenString: JWT
//
// Returns:
//   - *UserClaims: user claims
//   - error: error if the JWT verification failed
func (config *JWTConfig) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
