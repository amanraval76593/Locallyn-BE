package incident

import "context"

type Repository interface {
	GetNearbyIncidents(ctx context.Context, latitude float64, longitude float64, radius int) ([]Incident, error)
	GetIncident(ctx context.Context, id string) (*Incident, error)
	GetIncidentConfirmations(ctx context.Context, id string) ([]IncidentConfirmation, error)
	FindNearbyIncidentByCategory(ctx context.Context, latitude float64, longitude float64, radius int, category string) (*Incident, error)
	InsertIncident(ctx context.Context, incident *Incident) (*Incident, error)
}

type Service interface {
	GetNearbyIncidentService(ctx context.Context, req GetNearbyIncidentRequest) ([]Incident, error)
	GetIncidentService(ctx context.Context, req GetIncidentRequest) (*GetIncidentResponse, error)
	FindOrCreateIncidentService(ctx context.Context, latitude float64, longitude float64, radius int, category string) (*Incident, error)
	CreateIncidentService(ctx context.Context, latitude float64, longitude float64, radius int, category string) (*Incident, error)
}
