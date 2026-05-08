package feed

import (
	"context"
	"fmt"
	"locallyn-be/internal/common/constants"
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"
	"locallyn-be/pkg/database"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int, cursor *incidentCursor, limit int) ([]nearbyIncident, error) {
	query := `
		WITH nearby_incidents AS (
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
				ST_AsText(location) as location,
				ST_Distance(
					location,
					ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
				) as distance
			FROM incidents
			WHERE ST_DWithin(
				location,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
				$3
			)
			AND (expires_at IS NULL OR expires_at > NOW())
		),
		scored_incidents AS (
			SELECT
				*,
				(
					0.40 * COALESCE(trust_score, 0.5)
					+ 0.25 * GREATEST(0, 1 - (EXTRACT(EPOCH FROM (NOW() - created_at)) / 604800.0))
					+ 0.20 * LEAST(1, LN(1 + post_count + confirmation_count) / LN(101))
					+ 0.15 * GREATEST(0, 1 - (distance / NULLIF($3::FLOAT, 0)))
				) AS score
			FROM nearby_incidents
		)
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
			location,
			score
		FROM scored_incidents
	`

	args := []any{longitude, latitude, radius}
	if cursor != nil {
		query += fmt.Sprintf(`
			WHERE score < $%d
			OR (score = $%d AND created_at < $%d)
			OR (score = $%d AND created_at = $%d AND id::text < $%d)
		`, len(args)+1, len(args)+1, len(args)+2, len(args)+1, len(args)+2, len(args)+3)
		args = append(args, cursor.Score, cursor.CreatedAt, cursor.ID)
	}

	query += ` ORDER BY score DESC, created_at DESC, id::text DESC`

	if limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
		args = append(args, limit)
	}

	rows, err := database.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incidents := make([]nearbyIncident, 0)
	for rows.Next() {
		var item nearbyIncident

		if err := rows.Scan(
			&item.Incident.ID,
			&item.Incident.Title,
			&item.Incident.Category,
			&item.Incident.PostCount,
			&item.Incident.ConfirmationCount,
			&item.Incident.TrustScore,
			&item.Incident.CreatedAt,
			&item.Incident.UpdatedAt,
			&item.Incident.ExpiresAt,
			&item.Incident.Location,
			&item.Score,
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

func (r *repository) GetIncidentPosts(ctx context.Context, incidentID string, cursor *postCursor, limit int) ([]rankedPost, error) {
	query := `
		WITH scored_posts AS (
			SELECT
				p.id,
				p.user_id,
				p.incident_id,
				p.content,
				ST_AsText(p.location) as location,
				p.radius,
				p.identity_type,
				p.post_type,
				p.trust_score,
				p.media_urls,
				p.created_at,
				p.expires_at,
				p.is_deleted,
				p.is_flagged,
				(
					0.45 * COALESCE(p.trust_score, 0.5)
					+ 0.30 * GREATEST(0, 1 - (EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 604800.0))
					+ 0.15 * LEAST(1, LN(1 + COALESCE(feedback_counts.feedback_count, 0) + COALESCE(report_counts.report_count, 0)) / LN(51))
					+ 0.10 * COALESCE(
						GREATEST(
							0,
							1 - (
								ST_Distance(p.location, i.location)
								/ NULLIF(p.radius::FLOAT, 0)
							)
						),
						0
					)
				) AS score
			FROM posts p
			JOIN incidents i ON i.id = p.incident_id
			LEFT JOIN (
				SELECT post_id, COUNT(*)::FLOAT AS feedback_count
				FROM post_feedback
				GROUP BY post_id
			) feedback_counts ON feedback_counts.post_id = p.id
			LEFT JOIN (
				SELECT post_id, COUNT(*)::FLOAT AS report_count
				FROM post_reports
				GROUP BY post_id
			) report_counts ON report_counts.post_id = p.id
			WHERE p.incident_id = $1
			AND p.is_deleted = FALSE
			AND p.is_flagged = FALSE
			AND (p.expires_at IS NULL OR p.expires_at > NOW())
		)
		SELECT
			id,
			user_id,
			incident_id,
			content,
			location,
			radius,
			identity_type,
			post_type,
			trust_score,
			media_urls,
			created_at,
			expires_at,
			is_deleted,
			is_flagged,
			score
		FROM scored_posts
	`

	args := []any{incidentID}
	if cursor != nil {
		query += fmt.Sprintf(`
			WHERE score < $%d
			OR (score = $%d AND created_at < $%d)
			OR (score = $%d AND created_at = $%d AND id::text < $%d)
		`, len(args)+1, len(args)+1, len(args)+2, len(args)+1, len(args)+2, len(args)+3)
		args = append(args, cursor.Score, cursor.CreatedAt, cursor.ID)
	}

	query += ` ORDER BY score DESC, created_at DESC, id::text DESC`

	if limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
		args = append(args, limit)
	}

	rows, err := database.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]rankedPost, 0)
	for rows.Next() {
		var item rankedPost

		if err := rows.Scan(
			&item.Post.ID,
			&item.Post.UserID,
			&item.Post.IncidentID,
			&item.Post.Content,
			&item.Post.Location,
			&item.Post.Radius,
			&item.Post.IdentityType,
			&item.Post.PostType,
			&item.Post.TrustScore,
			&item.Post.MediaURLs,
			&item.Post.CreatedAt,
			&item.Post.ExpiresAt,
			&item.Post.IsDeleted,
			&item.Post.IsFlagged,
			&item.Score,
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

func (r *repository) GetNearbyBroadcastPosts(ctx context.Context, latitude float64, longitude float64, radius int, cursor *postCursor, limit int) ([]post.Post, error) {
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
	`

	args := []any{longitude, latitude, radius, constants.PostTypeBroadcast}
	if cursor != nil {
		query += fmt.Sprintf(`
			AND (
				created_at < $%d
				OR (created_at = $%d AND id < $%d)
			)
		`, len(args)+1, len(args)+1, len(args)+2)
		args = append(args, cursor.CreatedAt, cursor.ID)
	}

	query += ` ORDER BY created_at DESC, id DESC`

	if limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
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
