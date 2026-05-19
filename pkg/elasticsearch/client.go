package elasticsearch

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var Client *http.Client
var BaseURL string

func InitElasticsearch(elasticsearchURL string) {
	BaseURL = strings.TrimRight(elasticsearchURL, "/")
	Client = &http.Client{
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, BaseURL, nil)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch request: %v", err)
	}

	resp, err := Client.Do(req)
	if err != nil {
		log.Fatalf("Error initializing Elasticsearch: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Fatalf("Elasticsearch ping failed: %s: %s", resp.Status, string(body))
	}

	log.Printf("Elasticsearch client initialized: %s", BaseURL)
}
