package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jamcunha/expense-tracker/internal/utils"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { // Authorization: Bearer <token>
			token, err := validateJWT(strings.Split(r.Header.Get("Authorization"), " ")[1])
			if err != nil {
				utils.WriteJSON(
					w,
					http.StatusUnauthorized,
					utils.ApiError{Error: "Invalid Token"},
				)
				return
			}

			if !token.Valid {
				utils.WriteJSON(
					w,
					http.StatusUnauthorized,
					utils.ApiError{Error: "Invalid Token"},
				)
				return
			}

			userID, err := token.Claims.GetSubject()
			if err != nil {
				utils.WriteJSON(
					w,
					http.StatusUnauthorized,
					utils.ApiError{Error: "Invalid Token"},
				)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
}
