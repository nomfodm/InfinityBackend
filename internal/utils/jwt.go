package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenMalformed   = errors.New("provided string is not a token")
	ErrInvalidSignature = errors.New("invalid jwt signature")
	ErrTokenExpired     = errors.New("provided token expired")
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateJWTForUser(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	lifetime_minutes, _ := strconv.ParseUint(os.Getenv("JWT_LIFETIME"), 10, 32)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(lifetime_minutes) * time.Minute)),
		Subject:   fmt.Sprint(userID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	return signedToken, err
}

func ParseUserJWT(token string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	switch {
	case parsedToken.Valid:
		subject, _ := parsedToken.Claims.GetSubject()
		userID, _ := strconv.ParseUint(subject, 10, 32)
		return uint(userID), nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return 0, ErrTokenMalformed
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return 0, ErrInvalidSignature
	case errors.Is(err, jwt.ErrTokenExpired):
		return 0, ErrTokenExpired
	default:
		return 0, err
	}
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateRefreshToken() string {
	randomString := RandStringBytes(10)
	refreshToken, _ := HashStringToBcrypt(randomString)
	return refreshToken
}
