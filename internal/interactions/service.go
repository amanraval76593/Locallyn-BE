package interactions

import (
	"context"
	"errors"
	"locallyn-be/internal/common/auth"
	"locallyn-be/pkg/database"

	"github.com/google/uuid"
)

var (
	ErrIncidentNotFound    = errors.New("incident not found")
	ErrPostNotFound        = errors.New("post not found")
	ErrAlreadyConfirmed    = errors.New("incident already confirmed by user")
	ErrAlreadyGaveFeedback = errors.New("feedback already submitted for this post")
	ErrOwnPostFeedback     = errors.New("you cannot give feedback to your own post")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ConfirmIncidentService(ctx context.Context, req ConfirmIncidentRequest, claims *auth.Claims) error {
	if claims == nil || claims.UserId == "" {
		return errors.New("invalid authenticated user")
	}

	userID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	incidentId, err := uuid.Parse(req.IncidentId)
	if err != nil {
		return errors.New("invalid incident ID format")
	}

	err = database.WithTransaction(ctx, func(txCtx context.Context) error {

		err := s.repo.CreateConfirmIncident(txCtx, incidentId, userID)

		if err != nil {
			return err
		}

		err = s.repo.UpdateConfirmIncidentCount(txCtx, incidentId)

		if err != nil {
			return err
		}
		return nil
	})

	return err

}

func (s *service) PostFeedbackService(ctx context.Context, req PostFeedbackRequest, claims *auth.Claims) error {
	if claims == nil || claims.UserId == "" {
		return errors.New("invalid authenticated user")
	}

	userID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return errors.New("invalid post ID format")
	}

	ownerID, err := s.repo.GetPostOwnerID(ctx, postID)
	if err != nil {
		return err
	}

	if ownerID != nil && *ownerID == userID {
		return ErrOwnPostFeedback
	}

	return s.repo.CreatePostFeedback(ctx, postID, userID, req.Feedback)
}
