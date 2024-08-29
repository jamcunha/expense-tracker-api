package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { // Authorization: Bearer <token>
			if r.Header.Get("Authorization") == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "No Token Provided"}`))
				return
			}

			token, err := validateJWT(strings.Split(r.Header.Get("Authorization"), " ")[1])
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "Invalid Token"}`))
				return
			}

			if !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "Invalid Token"}`))
				return
			}

			userIDString, err := token.Claims.GetSubject()
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "Invalid Token"}`))
				return
			}

			userID, err := uuid.Parse(userIDString)
			if err != nil {
				fmt.Printf("Error parsing UUID: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "Invalid Token"}`))
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" { // TODO: find a better way to handle this
		fmt.Println("JWT_SECRET not set")
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
}
