package feed

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidCursor = errors.New("invalid cursor")

type postCursor struct {
	Score     float64   `json:"score,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

type incidentCursor struct {
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func encodeCursor(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodePostCursor(cursor string) (*postCursor, error) {
	if cursor == "" {
		return nil, nil
	}

	var value postCursor
	if err := decodeCursor(cursor, &value); err != nil {
		return nil, err
	}

	if value.CreatedAt.IsZero() || value.ID == "" {
		return nil, ErrInvalidCursor
	}

	return &value, nil
}

func decodeIncidentCursor(cursor string) (*incidentCursor, error) {
	if cursor == "" {
		return nil, nil
	}

	var value incidentCursor
	if err := decodeCursor(cursor, &value); err != nil {
		return nil, err
	}

	if value.CreatedAt.IsZero() || value.ID == "" {
		return nil, ErrInvalidCursor
	}

	return &value, nil
}

func decodeCursor(cursor string, value any) error {
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return ErrInvalidCursor
	}

	if err := json.Unmarshal(data, value); err != nil {
		return ErrInvalidCursor
	}

	return nil
}
