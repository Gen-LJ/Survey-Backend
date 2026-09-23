package jwt

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MinSecretLength is the shortest HS256 key worth accepting.
const MinSecretLength = 32

var ErrMissingSecret = errors.New("JWT_SECRET is not set")

// TokenTTL is how long an issued token stays valid.
var TokenTTL = 24 * time.Hour

// secret reads the signing key, refusing an empty one. Without this an unset
// JWT_SECRET would sign with a zero-length key, and every forged token would
// then validate.
func secret() ([]byte, error) {
	value := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if value == "" {
		return nil, ErrMissingSecret
	}

	return []byte(value), nil
}

// CheckSecret reports whether a usable signing key is configured. Call it at
// boot so a misconfigured deploy fails immediately instead of at first login.
func CheckSecret() error {
	_, err := secret()
	return err
}

func GenerateToken(email string) (string, error) {
	key, err := secret()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(TokenTTL).Unix(),
	})

	return token.SignedString(key)
}

func ValidateToken(tokenStr string) (string, error) {
	key, err := secret()
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
		return "", errors.New("email not found in token")
	}

	return "", errors.New("invalid token")
}
