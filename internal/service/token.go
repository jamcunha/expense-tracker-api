package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	DB               *pgx.Conn
	Queries          *repository.Queries
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExp     time.Duration
	JWTRefreshExp    time.Duration
}

func (s *Token) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := s.Queries.GetUserByEmail(r.Context(), body.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	} else if err != nil {
		fmt.Print("failed to query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !comparePassword(u.Password, body.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	}

	accessToken, err := createJWT(u, s.JWTAccessSecret, s.JWTAccessExp)
	if err != nil {
		fmt.Println("failed to create access token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken, err := createJWT(u, s.JWTRefreshSecret, s.JWTRefreshExp)
	if err != nil {
		fmt.Println("failed to create refresh token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := json.Marshal(struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (s *Token) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := validateJWT(body.RefreshToken, s.JWTRefreshSecret)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid token"}`))
		return
	}

	if !token.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid token"}`))
		return
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid token"}`))
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid token"}`))
	}

	u, err := s.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Print("failed to query:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	accessToken, err := createJWT(u, s.JWTAccessSecret, s.JWTAccessExp)
	if err != nil {
		fmt.Println("failed to create access token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := json.Marshal(struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	})
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func createJWT(
	user repository.User,
	jwtSecret string,
	jwtExpiration time.Duration,
) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtExpiration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "expense-tracker",
		Subject:   user.ID.String(),
	})

	return token.SignedString([]byte(jwtSecret))
}

// TODO: find a way to handle duplicate code (jwt internal package?)
//		 also validate may need more error handling

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

func comparePassword(userPassword string, givenPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword)) == nil
}
