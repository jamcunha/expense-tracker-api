package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jamcunha/expense-tracker/internal/repository"
	"github.com/jamcunha/expense-tracker/internal/service"
	"github.com/jamcunha/expense-tracker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService(t *testing.T) {
	var (
		mockDB      *mocks.MockDB
		mockQuerier *mocks.MockQuerier
		userService service.User
	)

	now := time.Now()
	expectedPassword := "passwd"
	encPassword, err := bcrypt.GenerateFromPassword([]byte(expectedPassword), bcrypt.DefaultCost)
	if !assert.Nilf(t, err, "expected to generate hash from password") {
		t.FailNow()
	}

	expectedUser := repository.User{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "test",
		Email:     "test@email.com",
		Password:  string(encPassword),
	}

	setup := func() {
		mockDB = &mocks.MockDB{}
		mockQuerier = &mocks.MockQuerier{}

		userService = service.User{
			DB:      mockDB,
			Queries: mockQuerier,
		}
	}

	t.Run("it should create a user", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("CreateUser", ctx, mock.Anything).Return(expectedUser, nil)

		user, err := userService.Create(
			ctx,
			expectedUser.Name,
			expectedUser.Email,
			expectedPassword,
		)

		assert.Nil(t, err)
		assertUserEqual(t, &expectedUser, &user)
	})

	t.Run("it should get a user by it's id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetUserByID", ctx, expectedUser.ID).Return(expectedUser, err)

		user, err := userService.GetByID(ctx, expectedUser.ID)

		assert.Nil(t, err)
		assertUserEqual(t, &expectedUser, &user)
	})

	t.Run("it should return an error if there is no user with the given id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetUserByID", ctx, mock.Anything).
			Return(repository.User{}, pgx.ErrNoRows)

		user, err := userService.GetByID(ctx, uuid.New())

		assert.ErrorIs(t, err, service.ErrUserNotFound)
		assert.Empty(t, user)
	})

	t.Run("it should delete a user given it's id", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("DeleteUser", ctx, expectedUser.ID).Return(expectedUser, nil)

		user, err := userService.DeleteByID(ctx, expectedUser.ID)

		assert.Nil(t, err)
		assertUserEqual(t, &expectedUser, &user)
	})

	t.Run(
		"it should return an error when there is no user with given id to delete",
		func(t *testing.T) {
			setup()

			ctx := context.Background()

			mockQuerier.On("DeleteUser", ctx, mock.Anything).
				Return(repository.User{}, pgx.ErrNoRows)

			user, err := userService.DeleteByID(ctx, uuid.New())

			assert.ErrorIs(t, err, service.ErrUserNotFound)
			assert.Empty(t, user)
		},
	)
}

func assertUserEqual(t *testing.T, expected *repository.User, actual *repository.User) {
	t.Helper()

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Password, actual.Password)
}
