package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal/database"
	"github.com/jamcunha/expense-tracker/internal/model"
	"github.com/jamcunha/expense-tracker/internal/utils"
)

func (s *Server) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)},
		)
		return
	}

	userID := r.Context().Value("userID")

	dbCategory, err := s.DB.CreateCategory(r.Context(), database.CreateCategoryParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		UserID:    userID.(uuid.UUID),
	})
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error creating category: %v", err)},
		)
		return
	}

	category := model.DatabaseCategoryToCategory(dbCategory)
	utils.WriteJSON(w, http.StatusCreated, category)
}

func (s *Server) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: fmt.Sprintf("Invalid JSON: %v", err)},
		)
		return
	}

	userID := r.Context().Value("userID")

	err = s.DB.DeleteCategory(r.Context(), database.DeleteCategoryParams{
		ID:     params.ID,
		UserID: userID.(uuid.UUID),
	})
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error deleting category: %v", err)},
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetUserCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	dbCategories, err := s.DB.GetUserCategories(r.Context(), userID)
	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: fmt.Sprintf("Error getting user categories: %v", err)},
		)
		return
	}

	categories := make([]model.Category, 0, len(dbCategories))
	for _, dbCategory := range dbCategories {
		categories = append(categories, model.DatabaseCategoryToCategory(dbCategory))
	}

	utils.WriteJSON(w, http.StatusOK, categories)
}
