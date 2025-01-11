package service

import (
	"context"
	"errors"
	"fmt"
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

func (s *Token) Create(
	ctx context.Context,
	email, password string,
) (accessToken, refreshToken string, err error) {
	u, err := s.Queries.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrUserNotFound
	} else if err != nil {
		return "", "", err
	}

	if !comparePassword(u.Password, password) {
		return "", "", ErrWrongCredentials
	}

	accessToken, err = createJWT(u, s.JWTAccessSecret, s.JWTAccessExp)
	if err != nil {
		fmt.Println("failed to create access token:", err)
		return "", "", err
	}

	refreshToken, err = createJWT(u, s.JWTRefreshSecret, s.JWTRefreshExp)
	if err != nil {
		fmt.Println("failed to create refresh token:", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Token) Refresh(ctx context.Context, refreshToken string) (string, error) {
	token, err := validateJWT(refreshToken, s.JWTRefreshSecret)
	if err != nil {
		return "", ErrInvalidToken
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v", err)
		return "", ErrInvalidToken
	}

	u, err := s.Queries.GetUserByID(ctx, userID)
	if err != nil {
		fmt.Print("failed to query:", err)
		return "", err
	}

	accessToken, err := createJWT(u, s.JWTAccessSecret, s.JWTAccessExp)
	if err != nil {
		fmt.Println("failed to create access token:", err)
		return "", err
	}

	return accessToken, nil
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
