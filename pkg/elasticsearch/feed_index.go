package elasticsearch

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

func BootstrapFeedPostsIndex(indexName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := indexExists(ctx, indexName)
	if err != nil {
		log.Fatalf("Error checking Elasticsearch index %q: %v", indexName, err)
	}

	if exists {
		log.Printf("Elasticsearch index already exists: %s", indexName)
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		BaseURL+"/"+indexName,
		bytes.NewBufferString(feedPostsMapping),
	)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch index request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch index %q: %v", indexName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Fatalf("Elasticsearch index creation failed for %q: %s: %s", indexName, resp.Status, string(body))
	}

	log.Printf("Elasticsearch index created: %s", indexName)
}

func indexExists(ctx context.Context, indexName string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, BaseURL+"/"+indexName, nil)
	if err != nil {
		return false, err
	}

	resp, err := Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return false, errors.New(resp.Status + ": " + string(body))
	}
}
