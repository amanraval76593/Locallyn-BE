package incident

import "context"

type Repository interface {
	GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int) ([]Incident, error)
	GetIncident(ctx context.Context, id string) (*Incident, error)
}

type Service interface {
	GetNearbyIncidentService(ctx context.Context, req GetNearbyIncidentRequest) ([]Incident, error)
	GetIncidentService(ctx context.Context, req GetIncidentRequest) (*Incident, error)
}
