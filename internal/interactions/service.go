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
	ErrOwnPostReport       = errors.New("you cannot report your own post")
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

	return database.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.CreatePostFeedback(txCtx, postID, userID, req.Feedback); err != nil {
			return err
		}

		if err := s.repo.UpdatePostTrustScore(txCtx, postID); err != nil {
			return err
		}

		return s.repo.UpdateUserTrustScore(txCtx, *ownerID)
	})
}

func (s *service) PostReportService(ctx context.Context, req PostReportRequest, claims *auth.Claims) (*PostReport, error) {
	if claims == nil || claims.UserId == "" {
		return nil, errors.New("invalid authenticated user")
	}

	userId, err := uuid.Parse(claims.UserId)

	if err != nil {
		return nil, errors.New("Invalid user ID")
	}

	postId, err := uuid.Parse(req.PostId)

	if err != nil {
		return nil, errors.New("Invalid post ID")
	}

	ownerID, err := s.repo.GetPostOwnerID(ctx, postId)
	if err != nil {
		return nil, err
	}

	if ownerID != nil && *ownerID == userId {
		return nil, ErrOwnPostReport
	}

	var postReport *PostReport
	err = database.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		postReport, err = s.repo.CreatePostReport(txCtx, postId, userId, req.Reason)
		if err != nil {
			return err
		}

		return s.repo.UpdateUserTrustScore(txCtx, *ownerID)
	})

	if err != nil {
		return nil, err
	}

	return postReport, nil

}
