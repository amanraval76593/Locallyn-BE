package post

import (
	"context"
	"errors"
	"fmt"
	"locallyn-be/internal/common/auth"
	"locallyn-be/internal/common/constants"
	"locallyn-be/internal/incident"
	"locallyn-be/pkg/database"
	"locallyn-be/pkg/elasticsearch"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo            Repository
	incidentService incident.Service
	feedSearchIndex string
}

func NewService(repo Repository, incidentService incident.Service, feedSearchIndex string) Service {
	return &service{
		repo:            repo,
		incidentService: incidentService,
		feedSearchIndex: feedSearchIndex,
	}
}

func (s *service) CreatePostService(ctx context.Context, claims *auth.Claims, req CreatePostRequest) (*Post, error) {

	if claims == nil || claims.UserId == "" {
		return nil, errors.New("invalid authenticated user")
	}

	userID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var newPost *Post
	var incidentID *uuid.UUID
	var incidentRecord *incident.Incident

	err = database.WithTransaction(ctx, func(txCtx context.Context) error {
		if req.Type == constants.PostTypeIncident {
			record, err := s.incidentService.FindOrCreateIncidentService(
				txCtx,
				req.Latitude,
				req.Longitude,
				2000,
				req.Category,
			)
			if err != nil {
				return err
			}

			incidentRecord = record
			incidentID = &record.ID
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

		createdPost, err := s.repo.CreatePost(txCtx, post)
		if err != nil {
			return err
		}

		if incidentID != nil {
			if err := s.repo.UpdateIncidentPostCount(txCtx, incidentID); err != nil {
				return err
			}
		}

		if err := s.repo.UpdateUserPostCount(txCtx, &userID); err != nil {
			return err
		}

		newPost = createdPost

		return nil
	})
	if err != nil {
		return nil, err
	}

	indexCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	doc, err := buildFeedPostDocument(*newPost, incidentRecord)
	if err != nil {
		return nil, err
	}

	if err := elasticsearch.IndexFeedPost(indexCtx, s.feedSearchIndex, doc); err != nil {
		return nil, err
	}

	return newPost, nil
}

func (s *service) FetchPostByIdService(ctx context.Context, req FetchPostByIdRequest) (*FetchPostByIdResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	post, err := s.repo.FetchPostById(ctx, &postID)

	if err != nil {
		return nil, err
	}

	feedbacks, err := s.repo.FetchPostFeedbacks(ctx, &postID)
	if err != nil {
		return nil, err
	}

	return &FetchPostByIdResponse{
		Post:      *post,
		Feedbacks: feedbacks,
	}, nil

}
