package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuth(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { // Authorization: Bearer <token>
			if r.Header.Get("Authorization") == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "No Token Provided"}`))
				return
			}

			token, err := validateJWT(
				strings.Split(r.Header.Get("Authorization"), " ")[1],
				jwtSecret,
			)
			if errors.Is(err, ErrExpiredToken) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				w.Write([]byte(`{"error": "Token Expired"}`))
				return
			} else if err != nil {
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

var ErrExpiredToken = fmt.Errorf("token is expired")

func validateJWT(tokenString string, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}

		if exp.Before(time.Now()) {
			return nil, ErrExpiredToken
		}
	}

	return token, nil
}
