package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func IndexFeedPost(ctx context.Context, indexName string, doc FeedPostDocument) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal feed post document: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		BaseURL+"/"+indexName+"/_doc/"+doc.ID,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create feed post index request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("index feed post document %q: %w", doc.ID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("index feed post document %q failed: %s: %s", doc.ID, resp.Status, string(respBody))
	}

	return nil
}
