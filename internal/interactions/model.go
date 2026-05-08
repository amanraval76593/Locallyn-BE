package interactions

import (
	"locallyn-be/internal/common/constants"
	"time"

	"github.com/google/uuid"
)

type PostFeedback struct {
	Id           uuid.UUID              `db:"id"`
	PostId       uuid.UUID              `db:"post_id"`
	UserId       uuid.UUID              `db:"user_id"`
	FeedbackType constants.FeedbackType `db:"feedback_type"`
	CreatedAt    time.Time              `db:"created_at"`
}

type PostReport struct {
	ID        uuid.UUID `db:"id"`
	PostId    uuid.UUID `db:"post_id"`
	UserId    uuid.UUID `db:"user_id"`
	Reason    string    `db:"reason"`
	CreatedAt time.Time `db:"created_at"`
}
