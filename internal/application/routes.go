package application

import (
	"fmt"
	"net/http"

	"github.com/jamcunha/expense-tracker/internal/handler"
	"github.com/jamcunha/expense-tracker/internal/middleware"
	"github.com/jamcunha/expense-tracker/internal/repository/category"
	"github.com/jamcunha/expense-tracker/internal/repository/user"
)

func (a *App) loadV1Routes(prefix string) {
	a.router = http.NewServeMux()

	// Health check
	a.router.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	if err := a.loadUserRoutes(prefix + "/users"); err != nil {
		fmt.Println("Error loading user routes:", err)
		return
	}

	if err := a.loadCategoryRoutes(prefix + "/categories"); err != nil {
		fmt.Println("Error loading category routes:", err)
		return
	}
}

func (a *App) loadUserRoutes(prefix string) error {
	if a.router == nil {
		return fmt.Errorf("router not initialized")
	}

	userHandler := &handler.User{
		Repo: &user.SqlcRepo{
			DB: a.db,
		},
	}

	a.router.HandleFunc("GET "+prefix+"/{id}", userHandler.GetByID)
	a.router.HandleFunc("POST "+prefix, userHandler.Create)
	a.router.HandleFunc("DELETE "+prefix+"/{id}", userHandler.DeleteByID)
	a.router.HandleFunc("POST /login", userHandler.Login) // does not use prefix

	return nil
}

func (a *App) loadCategoryRoutes(prefix string) error {
	if a.router == nil {
		return fmt.Errorf("router not initialized")
	}

	categoryHandler := &handler.Category{
		Repo: &category.SqlcRepo{
			DB: a.db,
		},
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f) }

	a.router.Handle("GET "+prefix, jwtMiddleware(categoryHandler.GetAll))
	a.router.Handle("GET "+prefix+"/{id}", jwtMiddleware(categoryHandler.GetByID))
	a.router.Handle("POST "+prefix, jwtMiddleware(categoryHandler.Create))
	a.router.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(categoryHandler.DeleteByID))

	return nil
}
