package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("L*%F9wV285W$wo8Zy8c*KPBoK"))
var pass = os.Getenv("TODO_PASSWORD")

type Credentials struct {
	Password string `json:"password"`
}

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Invalid JSON"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" || creds.Password != expectedPassword {
		writeJson(w, http.StatusForbidden, map[string]string{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	hash := generateHash(expectedPassword)
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Hash: hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{"token": tokenString})
}

func generateHash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

// Middleware for authentication
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pass == "" {
			// If no password is set, allow access
			next(w, r)
			return
		}

		// Extract token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Compare token hash with current password
		if claims.Hash != generateHash(pass) {
			http.Error(w, "Token no longer valid", http.StatusUnauthorized)
			return
		}

		// Passed authentication
		next(w, r)
	})
}
