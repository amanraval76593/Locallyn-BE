package elasticsearch

import "time"

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type FeedPostDocument struct {
	ID                 string     `json:"id"`
	UserID             *string    `json:"user_id,omitempty"`
	IncidentID         *string    `json:"incident_id,omitempty"`
	IncidentTitle      *string    `json:"incident_title,omitempty"`
	IncidentCategory   *string    `json:"incident_category,omitempty"`
	IncidentTrustScore *float64   `json:"incident_trust_score,omitempty"`
	Content            string     `json:"content"`
	Location           GeoPoint   `json:"location"`
	Radius             int        `json:"radius"`
	IdentityType       string     `json:"identity_type"`
	PostType           string     `json:"post_type"`
	TrustScore         *float64   `json:"trust_score,omitempty"`
	MediaURLs          []string   `json:"media_urls,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
	IsDeleted          bool       `json:"is_deleted"`
	IsFlagged          bool       `json:"is_flagged"`
}
