package service

import (
	"errors"
)

var (
	ErrUserNotFound     = errors.New("User not found")
	ErrCategoryNotFound = errors.New("Category not found")
	ErrExpenseNotFound  = errors.New("Expense not found")
	ErrBudgetNotFound   = errors.New("Budget not found")

	ErrWrongCredentials = errors.New("Wrong Credentials")
	ErrExpiredToken     = errors.New("Token is expired")
	ErrInvalidToken     = errors.New("Token is invalid")
	ErrDecodeCursor     = errors.New("Error decoding page cursor")
)
