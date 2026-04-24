package interactions

import (
	"context"
	"errors"
	"locallyn-be/internal/common/constants"
	"locallyn-be/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateConfirmIncident(ctx context.Context, incidentId uuid.UUID, userId uuid.UUID) error {
	query := `INSERT INTO incident_confirmations(
		incident_id,
		user_id
		)
		VALUES(
		$1,
		$2
		)
	`
	_, err := database.Conn(ctx).Exec(ctx, query, incidentId.String(), userId.String())

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return ErrIncidentNotFound
			case "23505":
				return ErrAlreadyConfirmed
			}
		}
		return err
	}

	return nil

}

func (r *repository) UpdateConfirmIncidentCount(ctx context.Context, incidentId uuid.UUID) error {
	query := `
		UPDATE incidents
		SET confirmation_count=confirmation_count+1
		WHERE id =$1
	`

	_, err := database.Conn(ctx).Exec(ctx, query, incidentId.String())

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetPostOwnerID(ctx context.Context, postId uuid.UUID) (*uuid.UUID, error) {
	query := `
		SELECT user_id
		FROM posts
		WHERE id = $1
		AND is_deleted = FALSE
	`

	var ownerID uuid.UUID

	err := database.Conn(ctx).QueryRow(ctx, query, postId.String()).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &ownerID, nil
}

func (r *repository) CreatePostFeedback(ctx context.Context, postId uuid.UUID, userId uuid.UUID, feedbackType constants.FeedbackType) error {
	query := `
		INSERT INTO post_feedback (
			post_id,
			user_id,
			feedback_type
		)
		VALUES ($1, $2, $3)
	`

	_, err := database.Conn(ctx).Exec(ctx, query, postId.String(), userId.String(), feedbackType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyGaveFeedback
		}
		return err
	}

	return nil
}
