package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jamcunha/expense-tracker/internal"
	"github.com/jamcunha/expense-tracker/internal/service"
)

type Token struct {
	service service.Token
}

type JWTParams struct {
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExp     time.Duration
	JWTRefreshExp    time.Duration
}

func NewToken(db internal.DBConn, queries internal.Querier, jwtParams JWTParams) *Token {
	return &Token{
		service: service.Token{
			DB:               db,
			Queries:          queries,
			JWTAccessSecret:  jwtParams.JWTAccessSecret,
			JWTRefreshSecret: jwtParams.JWTRefreshSecret,
			JWTAccessExp:     jwtParams.JWTAccessExp,
			JWTRefreshExp:    jwtParams.JWTRefreshExp,
		},
	}
}

func (h *Token) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.service.Create(r.Context(), body.Email, body.Password)
	if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrWrongCredentials) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	} else if err != nil {
		fmt.Print("failed to query:", err)
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

func (h *Token) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_tokend"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, err := h.service.Refresh(r.Context(), body.RefreshToken)
	if errors.Is(err, service.ErrInvalidToken) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		w.Write([]byte(`{"error": "Invalid token"}`))
		return
	} else if err != nil {
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
