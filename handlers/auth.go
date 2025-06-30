package handlers

import (
	"errors"
	"sync"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaim struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var tokenBlacklist = make(map[string]time.Time)
var blacklistMutex = &sync.RWMutex{} 

func GenerateJWT(userID uuid.UUID) (string, error) {
	expiration_time := time.Now().Add(72 * time.Hour)
	claims := &JWTClaim{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration_time),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(signedToken string) (uuid.UUID, error) {
	revoked := IsJWTRevoked(signedToken)
	if revoked {
		return uuid.UUID{0}, errors.New("token has been revoked")
	}

	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Println("Token not able to parse")
		return uuid.UUID{0}, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return uuid.UUID{0}, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return uuid.UUID{0}, errors.New("token already expired")
	} 

	return claims.UserId, nil
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	headerParts := strings.Split(authHeader, " ")
	if (len(headerParts) != 2 || headerParts[0] != "Bearer") {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return headerParts[1], nil
}

func RevokeJWT(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaim{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.New("invalid token claims")
	}
    blacklistMutex.Lock()
    defer blacklistMutex.Unlock()
    tokenBlacklist[tokenString] = claims.ExpiresAt.Time
	return nil
}

func IsJWTRevoked(tokenString string) bool {
    blacklistMutex.RLock()
    expiry, exists := tokenBlacklist[tokenString]
    blacklistMutex.RUnlock()
    
    if !exists {
        return false
    }
    
    if time.Now().After(expiry) {
        blacklistMutex.Lock()
        delete(tokenBlacklist, tokenString)
        blacklistMutex.Unlock()
        return false
    }
    
    return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
} 

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}