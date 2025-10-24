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

	// 1) Build query
	esQuery := buildSearchQuery(listStarterQuery)

	// 2) Chuẩn bị body: nil nếu query = nil; hoặc match_all nếu bạn muốn luôn có object
	var body io.Reader
	if esQuery != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
			return nil, 0, fmt.Errorf("error encoding query: %w", err)
		}
		// nếu buf.Len()==0 → coi như nil
		if buf.Len() > 0 {
			body = &buf
		}
	}

	if esQuery == nil {
		return nil, 0, fmt.Errorf("Fail to search from ES with provided query")
	}

	// 3) Pagination phòng thủ
	page, limit := 1, 10
	if p := listStarterQuery.Pagination.Page; p != nil && *p > 0 {
		page = *p
	}
	if l := listStarterQuery.Pagination.Limit; l != nil && *l > 0 {
		limit = *l
	}
	from := (page - 1) * limit

	// 4) Gọi ES
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(starterIndexName),
		r.client.Search.WithBody(body), // <- có thể nil
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithFrom(from),
		r.client.Search.WithSize(limit),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf("elasticsearch error: %s", string(b))
	}

	// 5) Parse response có kiểu rõ ràng
	type hitSrc struct {
		ID int64 `json:"id"`
		// thêm các field khác nếu cần
	}
	type esResp struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source hitSrc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var out esResp
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, 0, fmt.Errorf("error parsing response: %w", err)
	}

	ids := make([]int64, 0, len(out.Hits.Hits))
	for _, h := range out.Hits.Hits {
		ids = append(ids, h.Source.ID)
	}
	return ids, out.Hits.Total.Value, nil
}

func (r *ElasticsearchStarterRepository) IndexStarter(ctx context.Context, starter *model.StarterESDoc) error {
	doc := r.toDocument(starter)
	// Convert to JSON
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

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
		_ = Body.Close()
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}

func (r *ElasticsearchStarterRepository) DeleteFromIndex(ctx context.Context, domain string) error {
	// Build delete by query request
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"domain.keyword": domain,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("error marshaling delete query: %w", err)
	}

	refresh := true
	req := esapi.DeleteByQueryRequest{
		Index:   []string{starterIndexName},
		Body:    bytes.NewReader(body),
		Refresh: &refresh,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

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
		_ = Body.Close()
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
