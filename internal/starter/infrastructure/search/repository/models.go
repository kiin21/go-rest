package repository

import "time"

// indexName is the Elasticsearch index name for starters
const indexName = "starters"

// StarterDocument represents the Elasticsearch document structure for starters
// This is optimized for search, different from MySQL model
type StarterDocument struct {
	ID            int64  `json:"id"`
	Domain        string `json:"domain"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Mobile        string `json:"mobile"`
	WorkPhone     string `json:"work_phone"`
	JobTitle      string `json:"job_title"`
	DepartmentID  *int64 `json:"department_id,omitempty"`
	LineManagerID *int64 `json:"line_manager_id,omitempty"`

	// Additional fields for search optimization
	FullText     string   `json:"full_text"`     // Combined text for full-text search
	SearchTokens []string `json:"search_tokens"` // Tokenized fields for better matching

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IndexedAt time.Time `json:"indexed_at"` // When indexed to ES
}

// IndexMappingJSON returns the Elasticsearch index mapping
const IndexMappingJSON = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "starter_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding", "starter_edge_ngram"]
        },
        "starter_search_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding"]
        }
      },
      "filter": {
        "starter_edge_ngram": {
          "type": "edge_ngram",
          "min_gram": 2,
          "max_gram": 10
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },
      "domain": {
        "type": "text",
        "analyzer": "starter_analyzer",
        "search_analyzer": "starter_search_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "email": {
        "type": "text",
        "analyzer": "starter_analyzer",
        "search_analyzer": "starter_search_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "mobile": {
        "type": "text",
        "analyzer": "starter_analyzer",
        "search_analyzer": "starter_search_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "work_phone": {
        "type": "text"
      },
      "job_title": {
        "type": "text",
        "analyzer": "starter_analyzer",
        "search_analyzer": "starter_search_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "department_id": {
        "type": "long"
      },
      "line_manager_id": {
        "type": "long"
      },
      "full_text": {
        "type": "text",
        "analyzer": "starter_analyzer",
        "search_analyzer": "starter_search_analyzer"
      },
      "search_tokens": {
        "type": "keyword"
      },
      "created_at": {
        "type": "date"
      },
      "updated_at": {
        "type": "date"
      },
      "indexed_at": {
        "type": "date"
      }
    }
  }
}
`
