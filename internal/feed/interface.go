package feed

import (
	"context"
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"
)

type Repository interface {
	GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int) ([]incident.Incident, error)
	GetIncidentPosts(ctx context.Context, incidentID string, limit int) ([]post.Post, error)
	GetNearbyBroadcastPosts(ctx context.Context, latitude float64, longitude float64, radius int) ([]post.Post, error)
	GetIncidentByID(ctx context.Context, incidentID string) (*incident.Incident, error)
}

type Service interface {
	GetFeedByLocationService(ctx context.Context, req GetFeedByLocationRequest) (*GetFeedByLocationResponse, error)
	GetIncidentPostsService(ctx context.Context, req GetIncidentPostsRequest) (*GetIncidentPostsResponse, error)
}
