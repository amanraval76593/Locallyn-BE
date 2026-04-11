package incident

import (
	"context"
	"locallyn-be/pkg/database"
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

	rows, err := database.DB.Query(ctx, query, longitude, latitude, radius)

	if err != nil {
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

	err := database.DB.QueryRow(ctx, query, id).Scan(
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
