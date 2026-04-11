package incident

import "context"

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetNearbyIncidentService(ctx context.Context, req GetNearbyIncidentRequest) ([]Incident, error) {

	incidents, err := s.repo.GetNearbyIncidents(ctx, req.Latitude, req.Longitude, req.Radius)

	if err != nil {
		return nil, err
	}

	return incidents, nil
}

func (s *service) GetIncidentService(ctx context.Context, req GetIncidentRequest) (*Incident, error) {
	incident, err := s.repo.GetIncident(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return incident, nil
}
