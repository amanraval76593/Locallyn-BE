package post

import (
	"context"
	"locallyn-be/internal/common/auth"

	"github.com/google/uuid"
)

type Repository interface {
	CreatePost(ctx context.Context, post *Post) (*Post, error)
	UpdateUserPostCount(ctx context.Context, userId *uuid.UUID) error
	UpdateIncidentPostCount(ctx context.Context, incidentId *uuid.UUID) error
}

type Service interface {
	CreatePostService(ctx context.Context, claims *auth.Claims, req CreatePostRequest) (*Post, error)
}
