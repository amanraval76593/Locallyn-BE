package incident

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID                uuid.UUID  `db:"id"`
	Location          string     `db:"location"`
	Title             string     `db:"title"`
	Category          string     `db:"category"`
	PostCount         int        `db:"post_count"`
	ConfirmationCount int        `db:"confirmation_count"`
	TrustScore        float64    `db:"trust_score"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	ExpiresAt         *time.Time `db:"expires_at"`
}

type IncidentConfirmation struct {
	ID         uuid.UUID `db:"id" json:"id"`
	IncidentID uuid.UUID `db:"incident_id" json:"incident_id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
