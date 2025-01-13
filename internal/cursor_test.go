package internal_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jamcunha/expense-tracker/internal"
	"github.com/stretchr/testify/assert"
)

func TestCursor(t *testing.T) {
	expTime := time.Now().UTC()
	expUUID := uuid.New()
	testCursor := internal.EncodeCursor(expTime, expUUID)

	t.Run("it should decode a valid cursor", func(t *testing.T) {
		gotTime, gotUUID, err := internal.DecodeCursor(testCursor)

		assert.Nil(t, err)

		assert.Equal(t, gotTime, expTime)
		assert.Equal(t, gotUUID, expUUID)
	})

	t.Run("it should throw an error with an invalid cursor", func(t *testing.T) {
		gotTime, gotUUID, err := internal.DecodeCursor("invalid cursor")

		assert.Error(t, err)

		assert.Zero(t, gotTime)
		assert.Zero(t, gotUUID)
	})
}
