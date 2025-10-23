package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	repo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type ElasticsearchStarterRepository struct {
	client *elasticsearch.Client
}

func NewElasticsearchStarterRepository(client *elasticsearch.Client) repo.StarterSearchRepository {
	return &ElasticsearchStarterRepository{client: client}
}

func (r *ElasticsearchStarterRepository) Search(
	ctx context.Context,
	listStarterQuery *starterquery.ListStartersQuery,
	buildSearchQuery repo.SearchQueryBuilder,
) ([]int64, int64, error) {
	// Build ES query using the provided builder function
	esQuery := buildSearchQuery(listStarterQuery)

	// Execute search
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, 0, fmt.Errorf("error encoding query: %w", err)
	}

	// Extract pagination from query (type assertion)

	page, limit := 1, 10

	if listStarterQuery.Pagination.Page != nil {
		page = *listStarterQuery.Pagination.Page
	}
	if listStarterQuery.Pagination.Limit != nil {
		limit = *listStarterQuery.Pagination.Limit
	}

	from := (page - 1) * limit

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(starterIndexName),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithFrom(from),
		r.client.Search.WithSize(limit),
	)

	if err != nil {
		return nil, 0, fmt.Errorf("error executing search: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO: Handle error like this (search for "if err != nil")
		}
	}(res.Body)

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf("elasticsearch error: %s", string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("error parsing response: %w", err)
	}

	// Extract hits
	hits := result["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	documents := hits["hits"].([]interface{})

	// Convert to domain entities
	starterIds := make([]int64, 0, len(documents))
	for _, doc := range documents {
		source := doc.(map[string]interface{})["_source"].(map[string]interface{})
		starterId, err := r.toStarterID(source)
		if err != nil {
			return nil, 0, err
		}
		starterIds = append(starterIds, starterId)
	}

	return starterIds, total, nil
}

func (r *ElasticsearchStarterRepository) IndexStarter(ctx context.Context, starter *model.StarterESDoc) error {
	doc := r.toDocument(starter)
	// Convert to JSON
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	// Index document
	req := esapi.IndexRequest{
		Index:      starterIndexName,
		DocumentID: fmt.Sprintf("%d", starter.ID()),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO: Handle error like this (search for "if err != nil")
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}

func (r *ElasticsearchStarterRepository) DeleteFromIndex(ctx context.Context, domain string) error {
	// TODO: implement
	return nil
}

// BulkIndex indexes multiple starters in bulk (for initial indexing or reindexing)
func (r *ElasticsearchStarterRepository) BulkIndex(ctx context.Context, starters []*model.StarterESDoc) error {
	if len(starters) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, starter := range starters {
		doc := r.toDocument(starter)

		// Bulk index format: action line + document line
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": starterIndexName,
				"_id":    fmt.Sprintf("%d", starter.ID()),
			},
		}

		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("error encoding meta: %w", err)
		}

		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return fmt.Errorf("error encoding document: %w", err)
		}
	}

	// Execute bulk request
	req := esapi.BulkRequest{
		Index: starterIndexName,
		Body:  &buf,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error executing bulk: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO: Handle error like this (search for "if err != nil")
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error in bulk response: %s", res.String())
	}

	return nil
}

// toDocument converts domain Starter to ES document
func (r *ElasticsearchStarterRepository) toDocument(starter *model.StarterESDoc) *StarterDocument {
	// Build full text for search
	fullText := strings.Join([]string{
		strconv.FormatInt(starter.ID(), 10),
		starter.Domain(),
		starter.Name(),
		starter.DepartmentName(),
		starter.BusinessUnitName(),
	}, " ")

	// TODO: diff between search tokens and full text
	// Build search tokens (for exact matching)
	tokens := []string{
		strconv.FormatInt(starter.ID(), 10),
		starter.Domain(),
		starter.Name(),
		starter.DepartmentName(),
		starter.BusinessUnitName(),
	}

	return NewStarterDocumentBuilder().
		ID(starter.ID()).
		Domain(starter.Domain()).
		Name(starter.Name()).
		DepartmentName(starter.DepartmentName()).
		BusinessUnitName(starter.BusinessUnitName()).
		FullText(fullText).
		SearchTokens(tokens).
		BuildPtr()
}

// toStarterID converts ES document to domain Starter
func (r *ElasticsearchStarterRepository) toStarterID(source map[string]interface{}) (int64, error) {
	id, err := extractInt64(source, "id")
	if err != nil {
		return 0, fmt.Errorf("failed to extract id: %w", err)
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid id: %d", id)
	}

	return id, nil
}

// Helper function để extract int64 từ map[string]interface{}
func extractInt64(m map[string]interface{}, key string) (int64, error) {
	value, exists := m[key]
	if !exists {
		return 0, fmt.Errorf("field %s not found", key)
	}

	if value == nil {
		return 0, fmt.Errorf("field %s is nil", key)
	}

	// Handle multiple numeric types
	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("field %s is not a number, got type: %T", key, value)
	}
}
