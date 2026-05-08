package feed

import (
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"
)

type GetFeedByLocationRequest struct {
	Latitude        float64 `form:"lat" binding:"required,min=-90,max=90"`
	Longitude       float64 `form:"long" binding:"required,min=-180,max=180"`
	Radius          int     `form:"rad" binding:"required,min=1,max=50"`
	Limit           int     `form:"limit" binding:"omitempty,min=1,max=50"`
	IncidentCursor  string  `form:"incident_cursor"`
	BroadcastCursor string  `form:"broadcast_cursor"`
}

type IncidentFeedItem struct {
	Incident incident.Incident `json:"incident"`
	Posts    []post.Post       `json:"posts"`
	Score    float64           `json:"score"`
}

type GetFeedByLocationResponse struct {
	Incidents           []IncidentFeedItem `json:"incidents"`
	Broadcasts          []post.Post        `json:"broadcasts"`
	NextIncidentCursor  *string            `json:"next_incident_cursor,omitempty"`
	NextBroadcastCursor *string            `json:"next_broadcast_cursor,omitempty"`
	HasMoreIncidents    bool               `json:"has_more_incidents"`
	HasMoreBroadcasts   bool               `json:"has_more_broadcasts"`
	Limit               int                `json:"limit"`
}

type GetIncidentPostsRequest struct {
	IncidentID string `uri:"incidentId" binding:"required,uuid"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=50"`
	Cursor     string `form:"cursor"`
}

type GetIncidentPostsResponse struct {
	Incident   incident.Incident `json:"incident"`
	Posts      []post.Post       `json:"posts"`
	NextCursor *string           `json:"next_cursor,omitempty"`
	HasMore    bool              `json:"has_more"`
	Limit      int               `json:"limit"`
}
