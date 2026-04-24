package feed

import (
	"context"
	"locallyn-be/internal/common/constants"
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"
	"locallyn-be/pkg/database"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int) ([]incident.Incident, error) {
	query := `
		SELECT
			id,
			title,
			category,
			post_count,
			confirmation_count,
			trust_score,
			created_at,
			updated_at,
			expires_at,
			ST_AsText(location) as location
		FROM incidents
		WHERE ST_DWithin(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3
		)
		AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY ST_Distance(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		)
	`

	rows, err := database.Conn(ctx).Query(ctx, query, longitude, latitude, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incidents := make([]incident.Incident, 0)
	for rows.Next() {
		var item incident.Incident

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Category,
			&item.PostCount,
			&item.ConfirmationCount,
			&item.TrustScore,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ExpiresAt,
			&item.Location,
		); err != nil {
			return nil, err
		}

		incidents = append(incidents, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

func (r *repository) GetIncidentPosts(ctx context.Context, incidentID string, limit int) ([]post.Post, error) {
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
		WHERE incident_id = $1
		AND is_deleted = FALSE
		AND is_flagged = FALSE
		AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC
	`

	args := []any{incidentID}
	if limit > 0 {
		query += ` LIMIT $2`
		args = append(args, limit)
	}

	rows, err := database.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]post.Post, 0)
	for rows.Next() {
		var item post.Post

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.IncidentID,
			&item.Content,
			&item.Location,
			&item.Radius,
			&item.IdentityType,
			&item.PostType,
			&item.TrustScore,
			&item.MediaURLs,
			&item.CreatedAt,
			&item.ExpiresAt,
			&item.IsDeleted,
			&item.IsFlagged,
		); err != nil {
			return nil, err
		}

		posts = append(posts, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *repository) GetNearbyBroadcastPosts(ctx context.Context, latitude float64, longitude float64, radius int) ([]post.Post, error) {
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
		WHERE post_type = $4
		AND incident_id IS NULL
		AND is_deleted = FALSE
		AND is_flagged = FALSE
		AND (expires_at IS NULL OR expires_at > NOW())
		AND ST_DWithin(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3
		)
		ORDER BY created_at DESC
	`

	rows, err := database.Conn(ctx).Query(ctx, query, longitude, latitude, radius, constants.PostTypeBroadcast)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]post.Post, 0)
	for rows.Next() {
		var item post.Post

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.IncidentID,
			&item.Content,
			&item.Location,
			&item.Radius,
			&item.IdentityType,
			&item.PostType,
			&item.TrustScore,
			&item.MediaURLs,
			&item.CreatedAt,
			&item.ExpiresAt,
			&item.IsDeleted,
			&item.IsFlagged,
		); err != nil {
			return nil, err
		}

		posts = append(posts, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *repository) GetIncidentByID(ctx context.Context, incidentID string) (*incident.Incident, error) {
	query := `
		SELECT
			id,
			title,
			category,
			post_count,
			confirmation_count,
			trust_score,
			created_at,
			updated_at,
			expires_at,
			ST_AsText(location) as location
		FROM incidents
		WHERE id = $1
		AND (expires_at IS NULL OR expires_at > NOW())
	`

	var item incident.Incident
	err := database.Conn(ctx).QueryRow(ctx, query, incidentID).Scan(
		&item.ID,
		&item.Title,
		&item.Category,
		&item.PostCount,
		&item.ConfirmationCount,
		&item.TrustScore,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.ExpiresAt,
		&item.Location,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
