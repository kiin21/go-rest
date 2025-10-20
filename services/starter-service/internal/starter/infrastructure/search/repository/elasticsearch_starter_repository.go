package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/kiin21/go-rest/pkg/httputil"
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
	query string,
	filter model.StarterListFilter,
	pg httputil.ReqPagination,
) ([]*model.Starter, int64, error) {

	esQuery := r.buildSearchQuery(query, filter)

	// Execute search
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, 0, fmt.Errorf("error encoding query: %w", err)
	}

	// Pagination
	from := (pg.Page - 1) * pg.Limit

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(starterIndexName),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithFrom(from),
		r.client.Search.WithSize(pg.Limit),
	)
	if err != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return nil, 0, fmt.Errorf("error executing search: %w", err)
	}

	if res.IsError() {
		return nil, 0, fmt.Errorf("elasticsearch error: %s", res.String())
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
	starters := make([]*model.Starter, 0, len(documents))
	for _, doc := range documents {
		source := doc.(map[string]interface{})["_source"].(map[string]interface{})
		starter, err := r.toDomain(source)
		if err != nil {
			return nil, 0, err
		}
		starters = append(starters, starter)
	}

	return starters, total, nil
}

func (r *ElasticsearchStarterRepository) IndexStarter(ctx context.Context, starter *model.Starter) error {
	doc := r.toDocument(starter)
	doc.IndexedAt = time.Now()

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
		Refresh:    "true", // Make immediately searchable (use "false" for production with high volume)
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (r *ElasticsearchStarterRepository) DeleteFromIndex(ctx context.Context, domain string) error {
	// First, search for the document by domain to get its ID
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"domain.keyword": domain,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return fmt.Errorf("error encoding query: %w", err)
	}

	// Delete by query
	req := esapi.DeleteByQueryRequest{
		Index: []string{starterIndexName},
		Body:  &buf,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	return nil
}

// BulkIndex indexes multiple starters in bulk (for initial indexing or reindexing)
func (r *ElasticsearchStarterRepository) BulkIndex(ctx context.Context, starters []*model.Starter) error {
	if len(starters) == 0 {
		return nil
	}

	var buf bytes.Buffer
	now := time.Now()

	for _, starter := range starters {
		doc := r.toDocument(starter)
		doc.IndexedAt = now

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
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in bulk response: %s", res.String())
	}

	return nil
}

// Helper methods

func (r *ElasticsearchStarterRepository) buildSearchQuery(
	query string, filter model.StarterListFilter,
) map[string]interface{} {
	var must []interface{}

	if query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
				"fields": []string{
					"domain^3",    // Boost domain matches (highest priority)
					"job_title^2", // Boost job title matches
					"email",
					"mobile",
					"full_text",
				},
				"type":           "best_fields",
				"fuzziness":      "AUTO", // Handle typos
				"prefix_length":  2,
				"max_expansions": 50,
			},
		})
	}

	// Apply filters
	if filter.DepartmentID != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"department_id": *filter.DepartmentID,
			},
		})
	}

	if filter.LineManagerID != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"line_manager_id": *filter.LineManagerID,
			},
		})
	}

	// Build final query
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": buildSearchSortClause(filter.SortBy, filter.SortOrder),
	}

	return esQuery
}

// toDocument converts domain Starter to Elasticsearch document
func (r *ElasticsearchStarterRepository) toDocument(starter *model.Starter) *StarterDocument {
	// Build full text for search
	fullText := strings.Join([]string{
		starter.Domain(),
		starter.Name(),
		starter.Email(),
		starter.Mobile(),
		starter.JobTitle(),
	}, " ")

	// Build search tokens (for exact matching)
	tokens := []string{
		starter.Domain(),
		starter.Name(),
		starter.Email(),
		starter.Mobile(),
	}

	return &StarterDocument{
		ID:            starter.ID(),
		Domain:        starter.Domain(),
		Name:          starter.Name(),
		Email:         starter.Email(),
		Mobile:        starter.Mobile(),
		WorkPhone:     starter.WorkPhone(),
		JobTitle:      starter.JobTitle(),
		DepartmentID:  starter.DepartmentID(),
		LineManagerID: starter.LineManagerID(),
		FullText:      fullText,
		SearchTokens:  tokens,
		CreatedAt:     starter.CreatedAt(),
		UpdatedAt:     starter.UpdatedAt(),
	}
}

// toDomain converts Elasticsearch document to domain Starter
func (r *ElasticsearchStarterRepository) toDomain(source map[string]interface{}) (*model.Starter, error) {
	id := int64(source["id"].(float64))
	domain := source["domain"].(string)

	// Handle name field - might be nil for old documents indexed before name field was added
	name := ""
	if val, ok := source["name"]; ok && val != nil {
		name = val.(string)
	}

	email := source["email"].(string)
	mobile := source["mobile"].(string)
	workPhone := source["work_phone"].(string)
	jobTitle := source["job_title"].(string)

	var departmentID, lineManagerID *int64
	if val, ok := source["department_id"]; ok && val != nil {
		dep := int64(val.(float64))
		departmentID = &dep
	}
	if val, ok := source["line_manager_id"]; ok && val != nil {
		mgr := int64(val.(float64))
		lineManagerID = &mgr
	}

	createdAt, _ := time.Parse(time.RFC3339, source["created_at"].(string))
	updatedAt, _ := time.Parse(time.RFC3339, source["updated_at"].(string))

	return model.Rehydrate(
		id,
		domain,
		name,
		email,
		mobile,
		workPhone,
		jobTitle,
		departmentID,
		lineManagerID,
		createdAt,
		updatedAt,
	)
}

func buildSearchSortClause(sortBy, sortOrder string) []interface{} {
	field := mapStarterSearchSortField(sortBy)
	order := strings.ToLower(sortOrder)
	if order != "desc" {
		order = "asc"
	}

	if field == "" {
		return []interface{}{
			map[string]interface{}{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
			map[string]interface{}{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		}
	}

	sortClause := []interface{}{
		map[string]interface{}{
			field: map[string]interface{}{
				"order": order,
			},
		},
	}

	if field != "_score" {
		sortClause = append(sortClause, map[string]interface{}{
			"_score": map[string]interface{}{
				"order": "desc",
			},
		})
	}

	return sortClause
}

func mapStarterSearchSortField(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "id":
		return "id"
	case "domain":
		return "domain.keyword"
	case "created_at":
		return "created_at"
	default:
		return ""
	}
}
