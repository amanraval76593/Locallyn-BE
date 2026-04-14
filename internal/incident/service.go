package incident

import (
	"context"
	"fmt"
)

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

func (s *service) FindOrCreateIncidentServixe(ctx context.Context, latitude float64, longitude float64, radius int, category string) (*Incident, error) {
	incident, err := s.repo.FindNearbyIncidentByCategory(ctx, latitude, longitude, radius, category)

	if err != nil {
		return s.CreateIncidentService(ctx, latitude, longitude, radius, category)
	}

	if incident != nil {
		return incident, nil
	}

	return s.CreateIncidentService(ctx, latitude, longitude, radius, category)
}

func (s *service) CreateIncidentService(ctx context.Context, latitude float64, longitude float64, radius int, category string) (*Incident, error) {
	incident := &Incident{
		Location: fmt.Sprintf("POINT(%f %f)", longitude, latitude),
		Title:    category,
		Category: category,
	}

	createdIncident, err := s.repo.InsertIncident(ctx, incident)
	if err != nil {
		return nil, err
	}

	return createdIncident, nil
}
