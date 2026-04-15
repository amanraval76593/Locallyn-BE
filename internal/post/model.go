package post

import (
	"locallyn-be/internal/common/constants"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID              `db:"id"`
	UserID       *uuid.UUID             `db:"user_id"`
	IncidentID   *uuid.UUID             `db:"incident_id"`
	Content      string                 `db:"content"`
	Location     string                 `db:"location"`
	Radius       int                    `db:"radius"`
	IdentityType constants.IdentityType `db:"identity_type"`
	PostType     constants.PostType     `db:"post_type"`
	TrustScore   *float64               `db:"trust_score"`
	MediaURLs    []string               `db:"media_urls"`
	CreatedAt    time.Time              `db:"created_at"`
	ExpiresAt    *time.Time             `db:"expires_at"`
	IsDeleted    bool                   `db:"is_deleted"`
	IsFlagged    bool                   `db:"is_flagged"`
}
