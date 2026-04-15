package post

import (
	"context"
	"errors"
	"fmt"
	"locallyn-be/internal/common/auth"
	"locallyn-be/internal/common/constants"
	"locallyn-be/internal/incident"

	"github.com/google/uuid"
)

type service struct {
	repo            Repository
	incidentService incident.Service
}

func NewService(repo Repository, incidentService incident.Service) Service {
	return &service{
		repo:            repo,
		incidentService: incidentService,
	}
}

func (s *service) CreatePostService(ctx context.Context, claims *auth.Claims, req CreatePostRequest) (*Post, error) {

	if claims == nil || claims.UserId == "" {
		return nil, errors.New("invalid authenticated user")
	}

	var incidentID *uuid.UUID

	if req.Type == constants.PostTypeIncident {
		incidentRecord, err := s.incidentService.FindOrCreateIncidentService(
			ctx,
			req.Latitude,
			req.Longitude,
			2000,
			req.Category,
		)
		if err != nil {
			return nil, err
		}

		incidentID = &incidentRecord.ID
	}

	userID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	post := &Post{
		UserID:       &userID,
		IncidentID:   incidentID,
		Content:      req.Content,
		Location:     fmt.Sprintf("POINT(%f %f)", req.Longitude, req.Latitude),
		Radius:       2000,
		IdentityType: req.Identity,
		PostType:     req.Type,
		MediaURLs:    req.MediaURLs,
	}

	newPost, err := s.repo.CreatePost(ctx, post)

	if err != nil {
		return nil, err
	}
	if incidentID != nil {
		err := s.repo.UpdateIncidentPostCount(ctx, incidentID)

		if err != nil {
			return nil, err
		}
	}
	err = s.repo.UpdateUserPostCount(ctx, &userID)

	if err != nil {
		return nil, err
	}
	return newPost, nil
}

func (s *service) FetchPostByIdService(ctx context.Context, req FetchPostByIdRequest) (*Post, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	post, err := s.repo.FetchPostById(ctx, &postID)

	if err != nil {
		return nil, err
	}

	return post, nil

}
