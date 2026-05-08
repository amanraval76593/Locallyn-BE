package incident

import (
	"context"
	"errors"
	"locallyn-be/pkg/database"

	"github.com/jackc/pgx/v5"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int) ([]Incident, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	incidents := make([]Incident, 0)

	for rows.Next() {
		var incident Incident

		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Category,
			&incident.PostCount,
			&incident.ConfirmationCount,
			&incident.TrustScore,
			&incident.CreatedAt,
			&incident.UpdatedAt,
			&incident.ExpiresAt,
			&incident.Location,
		)

		if err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil

}

func (r *repository) GetIncident(ctx context.Context, id string) (*Incident, error) {
	query := `SELECT 
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
		WHERE id=$1
		`

	var incident Incident

	err := database.Conn(ctx).QueryRow(ctx, query, id).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Category,
		&incident.PostCount,
		&incident.ConfirmationCount,
		&incident.TrustScore,
		&incident.CreatedAt,
		&incident.UpdatedAt,
		&incident.ExpiresAt,
		&incident.Location,
	)

	if err != nil {
		return nil, err
	}

	return &incident, nil
}

func (r *repository) GetIncidentConfirmations(ctx context.Context, id string) ([]IncidentConfirmation, error) {
	query := `
		SELECT
			id,
			incident_id,
			user_id,
			created_at
		FROM incident_confirmations
		WHERE incident_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.Conn(ctx).Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	confirmations := make([]IncidentConfirmation, 0)

	for rows.Next() {
		var confirmation IncidentConfirmation

		err := rows.Scan(
			&confirmation.ID,
			&confirmation.IncidentID,
			&confirmation.UserID,
			&confirmation.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		confirmations = append(confirmations, confirmation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return confirmations, nil
}

func (r *repository) FindNearbyIncidentByCategory(ctx context.Context, latitude float64, longitude float64, radius int, category string) (*Incident, error) {
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
		AND category = $4
		AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY ST_Distance(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		)
		LIMIT 1;
	`

	var incident Incident

	err := database.Conn(ctx).QueryRow(ctx, query, longitude, latitude, radius, category).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Category,
		&incident.PostCount,
		&incident.ConfirmationCount,
		&incident.TrustScore,
		&incident.CreatedAt,
		&incident.UpdatedAt,
		&incident.ExpiresAt,
		&incident.Location,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &incident, nil
}

func (r *repository) InsertIncident(ctx context.Context, incident *Incident) (*Incident, error) {
	query := `
		INSERT INTO incidents (
			location,
			title,
			category,
			expires_at
		)
		VALUES (
			ST_GeogFromText($1),
			$2,
			$3,
			$4
		)
		RETURNING
			id,
			ST_AsText(location) as location,
			title,
			category,
			post_count,
			confirmation_count,
			trust_score,
			created_at,
			updated_at,
			expires_at
	`

	var created Incident

	err := database.Conn(ctx).QueryRow(
		ctx,
		query,
		incident.Location,
		incident.Title,
		incident.Category,
		incident.ExpiresAt,
	).Scan(
		&created.ID,
		&created.Location,
		&created.Title,
		&created.Category,
		&created.PostCount,
		&created.ConfirmationCount,
		&created.TrustScore,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
