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

func (r *repository) UpdatePostTrustScore(ctx context.Context, postId uuid.UUID) error {
	query := `
		UPDATE posts
		SET trust_score = feedback_score.score
		FROM (
			SELECT
				COALESCE(
					COUNT(*) FILTER (WHERE feedback_type = $2)::FLOAT
					/ NULLIF(COUNT(*)::FLOAT, 0),
					0.5
				) AS score
			FROM post_feedback
			WHERE post_id = $1
		) AS feedback_score
		WHERE posts.id = $1
	`

	_, err := database.Conn(ctx).Exec(ctx, query, postId.String(), constants.FeedbackHelpful)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateUserTrustScore(ctx context.Context, userId uuid.UUID) error {
	query := `
		WITH feedback_counts AS (
			SELECT
				COUNT(*) FILTER (WHERE pf.feedback_type = $2)::FLOAT AS helpful_count,
				COUNT(*) FILTER (WHERE pf.feedback_type = $3)::FLOAT AS misleading_count
			FROM posts p
			JOIN post_feedback pf ON pf.post_id = p.id
			WHERE p.user_id = $1
		),
		report_counts AS (
			SELECT COUNT(*)::FLOAT AS report_count
			FROM posts p
			JOIN post_reports pr ON pr.post_id = p.id
			WHERE p.user_id = $1
		)
		UPDATE user_profiles
		SET
			trust_score = (
				(feedback_counts.helpful_count + 1.0)
				/
				(
					feedback_counts.helpful_count
					+ feedback_counts.misleading_count
					+ report_counts.report_count
					+ 2.0
				)
			),
			total_reports = report_counts.report_count::INT,
			updated_at = NOW()
		FROM feedback_counts, report_counts
		WHERE user_profiles.user_id = $1
	`

	_, err := database.Conn(ctx).Exec(
		ctx,
		query,
		userId.String(),
		constants.FeedbackHelpful,
		constants.FeedbackMisleading,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreatePostReport(ctx context.Context, postId uuid.UUID, userId uuid.UUID, reason string) (*PostReport, error) {
	query := `
		INSERT INTO post_reports(
			post_id,
			user_id,
			reason
		)
		VALUES(
		$1,
		$2,
		$3
		)
		RETURNING 
		id,
		post_id,
		user_id,
		reason,
		created_at
	`

	var postReport PostReport

	err := database.Conn(ctx).QueryRow(
		ctx,
		query,
		postId,
		userId,
		reason,
	).Scan(
		&postReport.ID,
		&postReport.PostId,
		&postReport.UserId,
		&postReport.Reason,
		&postReport.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &postReport, nil
}
