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

	err := database.DB.QueryRow(
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
	_, err := database.DB.Exec(ctx, query, userId.String())

	return err
}

func (r *repository) UpdateIncidentPostCount(ctx context.Context, incidentId *uuid.UUID) error {
	query := `
	UPDATE incidents
	SET post_count=post_count+1
	WHERE id=$1
	`

	_, err := database.DB.Exec(ctx, query, incidentId.String())

	return err
}
