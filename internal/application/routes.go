package application

import (
	"net/http"

	"github.com/jamcunha/expense-tracker/internal/middleware"
	"github.com/jamcunha/expense-tracker/internal/service"
)

func (a *App) loadRoutes(prefix string) {
	a.router = http.NewServeMux()

	// Health check
	a.router.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	r := http.NewServeMux()

	a.loadUserRoutes(r, "/users")
	a.loadTokenRoutes(r, "/token")
	a.loadCategoryRoutes(r, "/categories")
	a.loadExpenseRoutes(r, "/expenses")
	a.loadBudgetRoutes(r, "/budgets")

	a.router.Handle(prefix+"/", http.StripPrefix(prefix, r))
}

func (a *App) loadUserRoutes(r *http.ServeMux, prefix string) {
	userService := &service.User{
		DB:      a.DB,
		Queries: a.Queries,
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f, a.config.JWTAccessSecret) }

	// NOTE: to restore password, add a route that requests the email and sends a token to the user
	// and another route that receives the token and the new password
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(userService.GetByID))
	r.HandleFunc("POST "+prefix, userService.Create)
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(userService.DeleteByID))
}

func (a *App) loadTokenRoutes(r *http.ServeMux, prefix string) {
	tokenService := &service.Token{
		DB:               a.DB,
		Queries:          a.Queries,
		JWTAccessSecret:  a.config.JWTAccessSecret,
		JWTRefreshSecret: a.config.JWTRefreshSecret,
		JWTAccessExp:     a.config.JWTAccessExp,
		JWTRefreshExp:    a.config.JWTRefreshExp,
	}

	r.HandleFunc("POST "+prefix, tokenService.Create)
	r.HandleFunc("POST "+prefix+"/refresh", tokenService.Refresh)
}

func (a *App) loadCategoryRoutes(r *http.ServeMux, prefix string) {
	categoryService := &service.Category{
		DB:      a.DB,
		Queries: a.Queries,
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f, a.config.JWTAccessSecret) }

	r.Handle("GET "+prefix, jwtMiddleware(categoryService.GetAll))
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(categoryService.GetByID))
	r.Handle("POST "+prefix, jwtMiddleware(categoryService.Create))
	r.Handle("PUT "+prefix+"/{id}", jwtMiddleware(categoryService.Update))
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(categoryService.DeleteByID))
}

func (a *App) loadExpenseRoutes(r *http.ServeMux, prefix string) {
	expenseService := &service.Expense{
		DB:      a.DB,
		Queries: a.Queries,
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f, a.config.JWTAccessSecret) }

	r.Handle("GET "+prefix, jwtMiddleware(expenseService.GetAll))
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(expenseService.GetByID))
	r.Handle("GET "+prefix+"/category/{id}", jwtMiddleware(expenseService.GetByCategory))
	r.Handle("POST "+prefix, jwtMiddleware(expenseService.Create))
	r.Handle("PUT "+prefix+"/{id}", jwtMiddleware(expenseService.Update))
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(expenseService.DeleteByID))
}

func (a *App) loadBudgetRoutes(r *http.ServeMux, prefix string) {
	budgetService := &service.Budget{
		DB:      a.DB,
		Queries: a.Queries,
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f, a.config.JWTAccessSecret) }

	r.Handle("GET "+prefix, jwtMiddleware(budgetService.GetAll))
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(budgetService.GetByID))
	r.Handle("POST "+prefix, jwtMiddleware(budgetService.Create))
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(budgetService.DeleteByID))
}
