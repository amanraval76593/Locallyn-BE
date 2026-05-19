package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FeedPostSearchHit struct {
	ID     string
	Score  float64
	Cursor FeedPostSearchCursor
}

type FeedPostSearchCursor struct {
	Score     float64 `json:"score"`
	CreatedAt string  `json:"created_at"`
	ID        string  `json:"id"`
}

func SearchFeedPosts(ctx context.Context, indexName string, keyword string, latitude float64, longitude float64, radiusInMeters int, cursor *FeedPostSearchCursor, limit int) ([]FeedPostSearchHit, error) {
	requestBody := map[string]any{
		"size": limit,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query": keyword,
							"fields": []string{
								"content^3",
								"incident_title^2",
								"incident_category",
							},
						},
					},
				},
				"filter": []any{
					map[string]any{
						"term": map[string]any{
							"is_deleted": false,
						},
					},
					map[string]any{
						"term": map[string]any{
							"is_flagged": false,
						},
					},
					map[string]any{
						"geo_distance": map[string]any{
							"distance": fmt.Sprintf("%dm", radiusInMeters),
							"location": GeoPoint{
								Lat: latitude,
								Lon: longitude,
							},
						},
					},
				},
				"should": []any{
					map[string]any{
						"range": map[string]any{
							"expires_at": map[string]any{
								"gt": "now",
							},
						},
					},
					map[string]any{
						"bool": map[string]any{
							"must_not": map[string]any{
								"exists": map[string]any{
									"field": "expires_at",
								},
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
		"sort": []any{
			map[string]any{
				"_score": map[string]any{
					"order": "desc",
				},
			},
			map[string]any{
				"created_at": map[string]any{
					"order":  "desc",
					"format": "strict_date_optional_time_nanos",
				},
			},
			map[string]any{
				"id": map[string]any{
					"order": "desc",
				},
			},
		},
	}

	if cursor != nil {
		requestBody["search_after"] = []any{
			cursor.Score,
			cursor.CreatedAt,
			cursor.ID,
		}
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal feed search request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, BaseURL+"/"+indexName+"/_search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create feed search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search feed posts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("search feed posts failed: %s: %s", resp.Status, string(respBody))
	}

	var searchResponse struct {
		Hits struct {
			Hits []struct {
				ID     string            `json:"_id"`
				Score  float64           `json:"_score"`
				Source FeedPostDocument  `json:"_source"`
				Sort   []json.RawMessage `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("decode feed search response: %w", err)
	}

	hits := make([]FeedPostSearchHit, 0, len(searchResponse.Hits.Hits))
	for _, hit := range searchResponse.Hits.Hits {
		cursor, err := searchCursorFromHit(hit.Score, hit.Source.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"), hit.ID, hit.Sort)
		if err != nil {
			return nil, err
		}

		hits = append(hits, FeedPostSearchHit{
			ID:     hit.ID,
			Score:  hit.Score,
			Cursor: cursor,
		})
	}

	return hits, nil
}

func searchCursorFromHit(score float64, createdAt string, id string, sortValues []json.RawMessage) (FeedPostSearchCursor, error) {
	cursor := FeedPostSearchCursor{
		Score:     score,
		CreatedAt: createdAt,
		ID:        id,
	}

	if len(sortValues) == 0 {
		return cursor, nil
	}

	if len(sortValues) != 3 {
		return FeedPostSearchCursor{}, fmt.Errorf("unexpected feed search sort value count: %d", len(sortValues))
	}

	if err := json.Unmarshal(sortValues[0], &cursor.Score); err != nil {
		return FeedPostSearchCursor{}, fmt.Errorf("decode search cursor score: %w", err)
	}

	if err := json.Unmarshal(sortValues[1], &cursor.CreatedAt); err != nil {
		return FeedPostSearchCursor{}, fmt.Errorf("decode search cursor created_at: %w", err)
	}

	if err := json.Unmarshal(sortValues[2], &cursor.ID); err != nil {
		return FeedPostSearchCursor{}, fmt.Errorf("decode search cursor id: %w", err)
	}

	return cursor, nil
}
