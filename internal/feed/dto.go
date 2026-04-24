package feed

import (
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"
)

type GetFeedByLocationRequest struct {
	Latitude  float64 `form:"lat" binding:"required,min=-90,max=90"`
	Longitude float64 `form:"long" binding:"required,min=-180,max=180"`
	Radius    int     `form:"rad" binding:"required,min=1,max=50"`
}

type IncidentFeedItem struct {
	Incident incident.Incident `json:"incident"`
	Posts    []post.Post       `json:"posts"`
}

type GetFeedByLocationResponse struct {
	Incidents  []IncidentFeedItem `json:"incidents"`
	Broadcasts []post.Post        `json:"broadcasts"`
}

type GetIncidentPostsRequest struct {
	IncidentID string `uri:"incidentId" binding:"required,uuid"`
}

type GetIncidentPostsResponse struct {
	Incident incident.Incident `json:"incident"`
	Posts    []post.Post       `json:"posts"`
}
