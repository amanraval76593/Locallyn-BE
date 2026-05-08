package feed

import (
	"context"
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"
)

type nearbyIncident struct {
	Incident incident.Incident
	Score    float64
}

type rankedPost struct {
	Post  post.Post
	Score float64
}

type Repository interface {
	GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int, cursor *incidentCursor, limit int) ([]nearbyIncident, error)
	GetIncidentPosts(ctx context.Context, incidentID string, cursor *postCursor, limit int) ([]rankedPost, error)
	GetNearbyBroadcastPosts(ctx context.Context, latitude float64, longitude float64, radius int, cursor *postCursor, limit int) ([]post.Post, error)
	GetIncidentByID(ctx context.Context, incidentID string) (*incident.Incident, error)
}

type Service interface {
	GetFeedByLocationService(ctx context.Context, req GetFeedByLocationRequest) (*GetFeedByLocationResponse, error)
	GetIncidentPostsService(ctx context.Context, req GetIncidentPostsRequest) (*GetIncidentPostsResponse, error)
}
