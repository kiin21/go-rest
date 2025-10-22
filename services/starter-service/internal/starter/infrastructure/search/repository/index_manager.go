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

type IndexManager struct {
	client *elasticsearch.Client
}

func NewIndexManager(client *elasticsearch.Client) *IndexManager {
	return &IndexManager{client: client}
}

func (im *IndexManager) CreateIndex(ctx context.Context) error {
	req := esapi.IndicesExistsRequest{
		Index: []string{starterIndexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	createReq := esapi.IndicesCreateRequest{
		Index: starterIndexName,
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

func (im *IndexManager) DeleteIndex(ctx context.Context) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{starterIndexName},
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

func (im *IndexManager) ResetIndex(ctx context.Context) error {
	_ = im.DeleteIndex(ctx)

	return im.CreateIndex(ctx)
}

func (im *IndexManager) GetDocumentCount(ctx context.Context) (int64, error) {
	req := esapi.CountRequest{
		Index: []string{starterIndexName},
	}

	res, err := req.Do(ctx, im.client)
	if err != nil {
		return 0, fmt.Errorf("error counting documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
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

func (im *IndexManager) IsIndexEmpty(ctx context.Context) (bool, error) {
	count, err := im.GetDocumentCount(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
