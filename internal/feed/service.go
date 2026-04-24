package feed

import "context"

const metersPerKilometer = 1000

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetFeedByLocationService(ctx context.Context, req GetFeedByLocationRequest) (*GetFeedByLocationResponse, error) {
	radiusInMeters := req.Radius * metersPerKilometer

	incidents, err := s.repo.GetNearbyIncidents(ctx, req.Latitude, req.Longitude, radiusInMeters)
	if err != nil {
		return nil, err
	}

	incidentFeed := make([]IncidentFeedItem, 0, len(incidents))
	for _, incident := range incidents {
		posts, err := s.repo.GetIncidentPosts(ctx, incident.ID.String(), 2)
		if err != nil {
			return nil, err
		}

		incidentFeed = append(incidentFeed, IncidentFeedItem{
			Incident: incident,
			Posts:    posts,
		})
	}

	broadcasts, err := s.repo.GetNearbyBroadcastPosts(ctx, req.Latitude, req.Longitude, radiusInMeters)
	if err != nil {
		return nil, err
	}

	return &GetFeedByLocationResponse{
		Incidents:  incidentFeed,
		Broadcasts: broadcasts,
	}, nil
}

func (s *service) GetIncidentPostsService(ctx context.Context, req GetIncidentPostsRequest) (*GetIncidentPostsResponse, error) {
	incidentRecord, err := s.repo.GetIncidentByID(ctx, req.IncidentID)
	if err != nil {
		return nil, err
	}

	posts, err := s.repo.GetIncidentPosts(ctx, req.IncidentID, 0)
	if err != nil {
		return nil, err
	}

	return &GetIncidentPostsResponse{
		Incident: *incidentRecord,
		Posts:    posts,
	}, nil
}
