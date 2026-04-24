package post

import (
	"context"
	"locallyn-be/pkg/database"

	"github.com/google/uuid"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	radius := post.Radius
	if radius == 0 {
		radius = 2000
	}

	query := `
		INSERT INTO posts (
			user_id,
			incident_id,
			content,
			location,
			radius,
			identity_type,
			post_type,
			trust_score,
			media_urls
		)
		VALUES (
			$1,
			$2,
			$3,
			ST_GeogFromText($4),
			$5,
			$6,
			$7,
			$8,
			$9
		)
		RETURNING
			id,
			user_id,
			incident_id,
			content,
			ST_AsText(location) as location,
			radius,
			identity_type,
			post_type,
			trust_score,
			media_urls,
			created_at,
			expires_at,
			is_deleted,
			is_flagged
	`

	var created Post

	err := database.Conn(ctx).QueryRow(
		ctx,
		query,
		post.UserID,
		post.IncidentID,
		post.Content,
		post.Location,
		radius,
		post.IdentityType,
		post.PostType,
		post.TrustScore,
		post.MediaURLs,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.IncidentID,
		&created.Content,
		&created.Location,
		&created.Radius,
		&created.IdentityType,
		&created.PostType,
		&created.TrustScore,
		&created.MediaURLs,
		&created.CreatedAt,
		&created.ExpiresAt,
		&created.IsDeleted,
		&created.IsFlagged,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *repository) UpdateUserPostCount(ctx context.Context, userId *uuid.UUID) error {
	query := `
	UPDATE user_profiles
	SET total_posts=total_posts+1
	WHERE user_id=$1
	`
	_, err := database.Conn(ctx).Exec(ctx, query, userId.String())

	return err
}

func (r *repository) UpdateIncidentPostCount(ctx context.Context, incidentId *uuid.UUID) error {
	query := `
	UPDATE incidents
	SET post_count=post_count+1
	WHERE id=$1
	`

	_, err := database.Conn(ctx).Exec(ctx, query, incidentId.String())

	return err
}

func (r *repository) FetchPostById(ctx context.Context, postId *uuid.UUID) (*Post, error) {
	query := `
		SELECT
			id,
			user_id,
			incident_id,
			content,
			ST_AsText(location) as location,
			radius,
			identity_type,
			post_type,
			trust_score,
			media_urls,
			created_at,
			expires_at,
			is_deleted,
			is_flagged
		FROM posts
		WHERE id=$1
	`

	var post Post

	err := database.Conn(ctx).QueryRow(ctx, query, postId).Scan(
		&post.ID,
		&post.UserID,
		&post.IncidentID,
		&post.Content,
		&post.Location,
		&post.Radius,
		&post.IdentityType,
		&post.PostType,
		&post.TrustScore,
		&post.MediaURLs,
		&post.CreatedAt,
		&post.ExpiresAt,
		&post.IsDeleted,
		&post.IsFlagged,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *repository) FetchPostFeedbacks(ctx context.Context, postId *uuid.UUID) ([]PostFeedback, error) {
	query := `
		SELECT
			id,
			post_id,
			user_id,
			feedback_type,
			created_at
		FROM post_feedback
		WHERE post_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.Conn(ctx).Query(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := make([]PostFeedback, 0)

	for rows.Next() {
		var feedback PostFeedback

		err := rows.Scan(
			&feedback.ID,
			&feedback.PostID,
			&feedback.UserID,
			&feedback.FeedbackType,
			&feedback.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		feedbacks = append(feedbacks, feedback)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedbacks, nil
}
