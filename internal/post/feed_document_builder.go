package post

import (
	"fmt"
	"locallyn-be/internal/incident"
	"locallyn-be/pkg/elasticsearch"
)

func buildFeedPostDocument(postItem Post, incidentItem *incident.Incident) (elasticsearch.FeedPostDocument, error) {
	location, err := elasticsearch.ParseWKTPoint(postItem.Location)
	if err != nil {
		return elasticsearch.FeedPostDocument{}, fmt.Errorf("parse post location: %w", err)
	}

	doc := elasticsearch.FeedPostDocument{
		ID:           postItem.ID.String(),
		Content:      postItem.Content,
		Location:     location,
		Radius:       postItem.Radius,
		IdentityType: string(postItem.IdentityType),
		PostType:     string(postItem.PostType),
		TrustScore:   postItem.TrustScore,
		MediaURLs:    postItem.MediaURLs,
		CreatedAt:    postItem.CreatedAt,
		ExpiresAt:    postItem.ExpiresAt,
		IsDeleted:    postItem.IsDeleted,
		IsFlagged:    postItem.IsFlagged,
	}

	if postItem.UserID != nil {
		userID := postItem.UserID.String()
		doc.UserID = &userID
	}

	if postItem.IncidentID != nil {
		incidentID := postItem.IncidentID.String()
		doc.IncidentID = &incidentID
	}

	if incidentItem != nil {
		doc.IncidentTitle = &incidentItem.Title
		doc.IncidentCategory = &incidentItem.Category
		doc.IncidentTrustScore = &incidentItem.TrustScore
	}

	return doc, nil
}
