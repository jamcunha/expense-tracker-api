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

// TODO: Create a helper package and move JWT related logic to there
//       Not testing JWT logic here

func TestTokenService(t *testing.T) {
	var (
		mockDB       *mocks.MockDB
		mockQuerier  *mocks.MockQuerier
		tokenService service.Token
	)

	testPassword := "test"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if !assert.Nilf(t, err, "expected to generate hash from password") {
		t.FailNow()
	}

	now := time.Now()
	testUser := repository.User{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "Test",
		Email:     "test@email.com",
		Password:  string(hashedPassword),
	}

	setup := func() {
		mockDB = &mocks.MockDB{}
		mockQuerier = &mocks.MockQuerier{}

		tokenService = service.Token{
			DB:               mockDB,
			Queries:          mockQuerier,
			JWTAccessSecret:  "access-secret",
			JWTRefreshSecret: "refresh-secret",
			JWTAccessExp:     time.Duration(5 * time.Minute),
			JWTRefreshExp:    time.Duration(15 * time.Minute),
		}
	}

	t.Run("it should create tokens given correct credentials", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetUserByEmail", ctx, testUser.Email).Return(testUser, nil)

		acc, ref, err := tokenService.Create(
			ctx,
			testUser.Email,
			testPassword,
		)

		assert.Nil(t, err)
		assert.NotEmpty(t, acc)
		assert.NotEmpty(t, ref)
	})

	t.Run("it should return an error if email does not exist", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetUserByEmail", ctx, mock.Anything).
			Return(repository.User{}, pgx.ErrNoRows)

		acc, ref, err := tokenService.Create(
			ctx,
			testUser.Email,
			testPassword,
		)

		assert.ErrorIs(t, err, service.ErrUserNotFound)
		assert.Empty(t, acc)
		assert.Empty(t, ref)
	})

	t.Run("it should return an error if given wrong credentials", func(t *testing.T) {
		setup()

		ctx := context.Background()

		mockQuerier.On("GetUserByEmail", ctx, testUser.Email).Return(testUser, nil)

		acc, ref, err := tokenService.Create(ctx, testUser.Email, "wrong-password")

		assert.ErrorIs(t, err, service.ErrWrongCredentials)
		assert.Empty(t, acc)
		assert.Empty(t, ref)
	})

	// TODO: Test Refresh
	//       Now it's not possible because it uses private methods to validate
	//       the token. Return later when split JWT logic to a separate package
}
