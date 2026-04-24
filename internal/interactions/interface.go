package interactions

import (
	"context"
	"locallyn-be/internal/common/auth"
	"locallyn-be/internal/common/constants"

	"github.com/google/uuid"
)

type Repository interface {
	CreateConfirmIncident(ctx context.Context, incidentId uuid.UUID, userId uuid.UUID) error
	UpdateConfirmIncidentCount(ctx context.Context, incidentId uuid.UUID) error
	GetPostOwnerID(ctx context.Context, postId uuid.UUID) (*uuid.UUID, error)
	CreatePostFeedback(ctx context.Context, postId uuid.UUID, userId uuid.UUID, feedbackType constants.FeedbackType) error
}

type Service interface {
	ConfirmIncidentService(ctx context.Context, req ConfirmIncidentRequest, claims *auth.Claims) error
	PostFeedbackService(ctx context.Context, req PostFeedbackRequest, claims *auth.Claims) error
}
