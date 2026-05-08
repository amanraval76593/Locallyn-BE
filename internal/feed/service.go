package feed

import (
	"context"
	"locallyn-be/internal/post"
)

const metersPerKilometer = 1000
const defaultPageLimit = 20

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetFeedByLocationService(ctx context.Context, req GetFeedByLocationRequest) (*GetFeedByLocationResponse, error) {
	radiusInMeters := req.Radius * metersPerKilometer
	limit := normalizeLimit(req.Limit)

	decodedIncidentCursor, err := decodeIncidentCursor(req.IncidentCursor)
	if err != nil {
		return nil, err
	}

	broadcastCursor, err := decodePostCursor(req.BroadcastCursor)
	if err != nil {
		return nil, err
	}

	incidents, err := s.repo.GetNearbyIncidents(ctx, req.Latitude, req.Longitude, radiusInMeters, decodedIncidentCursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMoreIncidents := len(incidents) > limit
	if hasMoreIncidents {
		incidents = incidents[:limit]
	}

	incidentFeed := make([]IncidentFeedItem, 0, len(incidents))
	for _, incidentRecord := range incidents {
		rankedPosts, err := s.repo.GetIncidentPosts(ctx, incidentRecord.Incident.ID.String(), nil, 2)
		if err != nil {
			return nil, err
		}

		incidentFeed = append(incidentFeed, IncidentFeedItem{
			Incident: incidentRecord.Incident,
			Posts:    postsFromRankedPosts(rankedPosts),
			Score:    incidentRecord.Score,
		})
	}

	var nextIncidentCursor *string
	if hasMoreIncidents && len(incidents) > 0 {
		last := incidents[len(incidents)-1]
		cursor, err := encodeCursor(incidentCursor{
			Score:     last.Score,
			CreatedAt: last.Incident.CreatedAt,
			ID:        last.Incident.ID.String(),
		})
		if err != nil {
			return nil, err
		}
		nextIncidentCursor = &cursor
	}

	broadcasts, err := s.repo.GetNearbyBroadcastPosts(ctx, req.Latitude, req.Longitude, radiusInMeters, broadcastCursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMoreBroadcasts := len(broadcasts) > limit
	if hasMoreBroadcasts {
		broadcasts = broadcasts[:limit]
	}

	var nextBroadcastCursor *string
	if hasMoreBroadcasts && len(broadcasts) > 0 {
		last := broadcasts[len(broadcasts)-1]
		cursor, err := encodeCursor(postCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID.String(),
		})
		if err != nil {
			return nil, err
		}
		nextBroadcastCursor = &cursor
	}

	return &GetFeedByLocationResponse{
		Incidents:           incidentFeed,
		Broadcasts:          broadcasts,
		NextIncidentCursor:  nextIncidentCursor,
		NextBroadcastCursor: nextBroadcastCursor,
		HasMoreIncidents:    hasMoreIncidents,
		HasMoreBroadcasts:   hasMoreBroadcasts,
		Limit:               limit,
	}, nil
}

func (s *service) GetIncidentPostsService(ctx context.Context, req GetIncidentPostsRequest) (*GetIncidentPostsResponse, error) {
	incidentRecord, err := s.repo.GetIncidentByID(ctx, req.IncidentID)
	if err != nil {
		return nil, err
	}

	limit := normalizeLimit(req.Limit)

	cursor, err := decodePostCursor(req.Cursor)
	if err != nil {
		return nil, err
	}

	rankedPosts, err := s.repo.GetIncidentPosts(ctx, req.IncidentID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := len(rankedPosts) > limit
	if hasMore {
		rankedPosts = rankedPosts[:limit]
	}

	var nextCursor *string
	if hasMore && len(rankedPosts) > 0 {
		last := rankedPosts[len(rankedPosts)-1]
		cursor, err := encodeCursor(postCursor{
			Score:     last.Score,
			CreatedAt: last.Post.CreatedAt,
			ID:        last.Post.ID.String(),
		})
		if err != nil {
			return nil, err
		}
		nextCursor = &cursor
	}

	return &GetIncidentPostsResponse{
		Incident:   *incidentRecord,
		Posts:      postsFromRankedPosts(rankedPosts),
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      limit,
	}, nil
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultPageLimit
	}

	return limit
}

func postsFromRankedPosts(rankedPosts []rankedPost) []post.Post {
	posts := make([]post.Post, 0, len(rankedPosts))
	for _, item := range rankedPosts {
		posts = append(posts, item.Post)
	}

	return posts
}
