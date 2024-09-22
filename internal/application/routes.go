package application

import (
	"net/http"

	"github.com/jamcunha/expense-tracker/internal/handler"
	"github.com/jamcunha/expense-tracker/internal/middleware"
	"github.com/jamcunha/expense-tracker/internal/repository/budget"
	"github.com/jamcunha/expense-tracker/internal/repository/category"
	"github.com/jamcunha/expense-tracker/internal/repository/expense"
	"github.com/jamcunha/expense-tracker/internal/repository/user"
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
	a.loadCategoryRoutes(r, "/categories")
	a.loadExpenseRoutes(r, "/expenses")
	a.loadBudgetRoutes(r, "/budgets")

	a.router.Handle(prefix+"/", http.StripPrefix(prefix, r))
}

func (a *App) loadUserRoutes(r *http.ServeMux, prefix string) {
	userHandler := &handler.User{
		Repo: &user.SqlcRepo{
			DB:      a.DB,
			Queries: a.Queries,
		},
	}

	// NOTE: to restore password, add a route that requests the email and sends a token to the user
	// and another route that receives the token and the new password
	r.HandleFunc("GET "+prefix+"/{id}", userHandler.GetByID)
	r.HandleFunc("POST "+prefix, userHandler.Create)
	r.HandleFunc("DELETE "+prefix+"/{id}", userHandler.DeleteByID)
	r.HandleFunc("POST /login", userHandler.Login) // does not use prefix
}

func (a *App) loadCategoryRoutes(r *http.ServeMux, prefix string) {
	categoryHandler := &handler.Category{
		Repo: &category.SqlcRepo{
			DB:      a.DB,
			Queries: a.Queries,
		},
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f) }

	r.Handle("GET "+prefix, jwtMiddleware(categoryHandler.GetAll))
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(categoryHandler.GetByID))
	r.Handle("POST "+prefix, jwtMiddleware(categoryHandler.Create))
	r.Handle("PUT "+prefix+"/{id}", jwtMiddleware(categoryHandler.Update))
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(categoryHandler.DeleteByID))
}

func (a *App) loadExpenseRoutes(r *http.ServeMux, prefix string) {
	expenseHandler := &handler.Expense{
		Repo: &expense.SqlcRepo{
			DB:      a.DB,
			Queries: a.Queries,
		},
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f) }

	r.Handle("GET "+prefix, jwtMiddleware(expenseHandler.GetAll))
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(expenseHandler.GetByID))
	r.Handle("GET "+prefix+"/category/{id}", jwtMiddleware(expenseHandler.GetByCategory))
	r.Handle("POST "+prefix, jwtMiddleware(expenseHandler.Create))
	r.Handle("PUT "+prefix+"/{id}", jwtMiddleware(expenseHandler.Update))
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(expenseHandler.DeleteByID))
}

func (a *App) loadBudgetRoutes(r *http.ServeMux, prefix string) {
	budgetHandler := &handler.Budget{
		Repo: &budget.SqlcRepo{
			DB:      a.DB,
			Queries: a.Queries,
		},
	}

	jwtMiddleware := func(f http.HandlerFunc) http.Handler { return middleware.JWTAuth(f) }

	r.Handle("GET "+prefix, jwtMiddleware(budgetHandler.GetAll))
	r.Handle("GET "+prefix+"/{id}", jwtMiddleware(budgetHandler.GetByID))
	r.Handle("POST "+prefix, jwtMiddleware(budgetHandler.Create))
	r.Handle("DELETE "+prefix+"/{id}", jwtMiddleware(budgetHandler.DeleteByID))
}
