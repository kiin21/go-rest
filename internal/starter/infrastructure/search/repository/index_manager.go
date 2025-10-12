package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// IndexManager handles Elasticsearch index operations
type IndexManager struct {
	client *elasticsearch.Client
}

// NewIndexManager creates a new index manager
func NewIndexManager(client *elasticsearch.Client) *IndexManager {
	return &IndexManager{client: client}
}

// CreateIndex creates the starters index with proper mapping
func (im *IndexManager) CreateIndex(ctx context.Context) error {
	// Check if index exists
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	// Index already exists
	if res.StatusCode == 200 {
		return nil
	}

	// Create index with mapping
	createReq := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader([]byte(IndexMappingJSON)),
	}

	createRes, err := createReq.Do(ctx, im.client)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		body, _ := io.ReadAll(createRes.Body)
		return fmt.Errorf("error creating index: %s", string(body))
	}

	return nil
}

// DeleteIndex deletes the starters index
func (im *IndexManager) DeleteIndex(ctx context.Context) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return fmt.Errorf("error deleting index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error deleting index: %s", string(body))
	}

	return nil
}

// ResetIndex deletes and recreates the index
func (im *IndexManager) ResetIndex(ctx context.Context) error {
	// Delete index (ignore error if not exists)
	_ = im.DeleteIndex(ctx)

	// Create index
	return im.CreateIndex(ctx)
}

// GetIndexStats returns index statistics
func (im *IndexManager) GetIndexStats(ctx context.Context) (map[string]interface{}, error) {
	req := esapi.IndicesStatsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return nil, fmt.Errorf("error getting index stats: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error getting index stats: %s", string(body))
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("error parsing stats response: %w", err)
	}

	return stats, nil
}

// RefreshIndex refreshes the index to make all operations available for search
func (im *IndexManager) RefreshIndex(ctx context.Context) error {
	req := esapi.IndicesRefreshRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return fmt.Errorf("error refreshing index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error refreshing index: %s", string(body))
	}

	return nil
}

// GetDocumentCount returns the number of documents in the index
func (im *IndexManager) GetDocumentCount(ctx context.Context) (int64, error) {
	req := esapi.CountRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return 0, fmt.Errorf("error counting documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		// Index might not exist, return 0
		if res.StatusCode == 404 {
			return 0, nil
		}
		body, _ := io.ReadAll(res.Body)
		return 0, fmt.Errorf("error counting documents: %s", string(body))
	}

	var result struct {
		Count int64 `json:"count"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error parsing count response: %w", err)
	}

	return result.Count, nil
}

// IsIndexEmpty checks if the index exists and is empty
func (im *IndexManager) IsIndexEmpty(ctx context.Context) (bool, error) {
	count, err := im.GetDocumentCount(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
