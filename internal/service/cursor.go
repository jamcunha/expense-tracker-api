package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Util functions to encode and decode cursor for paged queries

func decodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, uuid.UUID{}, err
	}

	t, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}

	id, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}

	return t, id, nil
}

func encodeCursor(t time.Time, uuid uuid.UUID) string {
	return base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid.String()),
	))
}
