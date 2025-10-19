# Starter Management API

A Go-based REST API for managing employee onboarding with real-time search via Elasticsearch and event-driven sync using Kafka.

## Architecture Overview

```mermaid
flowchart TB
    subgraph Client Layer
        UI[Swagger UI / REST Client]
    end
    
    subgraph Application Layer
        API[Gin Router]
        AppSvc[Application Service]
        SearchSvc[Search Service]
    end
    
    subgraph Infrastructure
        Repo[GORM Repository]
        ESRepo[ES Repository]
        Producer[Kafka Producer]
        Consumer[Kafka Consumer]
    end
    
    subgraph Data Layer
        MySQL[(MySQL)]
        ES[(Elasticsearch)]
        Kafka[(Kafka)]
    end
    
    UI -->|HTTP| API
    API --> AppSvc
    AppSvc --> Repo --> MySQL
    AppSvc --> SearchSvc
    SearchSvc -.->|Query| ES
    SearchSvc -->|Events| Producer --> Kafka
    Kafka --> Consumer --> ESRepo --> ES
    
    style MySQL fill:#4479A1,color:#fff
    style ES fill:#005571,color:#fff
    style Kafka fill:#231F20,color:#fff
```

## Data Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant A as API Handler
    participant S as Service
    participant M as MySQL
    participant K as Kafka
    participant E as Elasticsearch

    Note over C,E: Create/Update Flow
    C->>A: POST /starters
    A->>S: CreateStarter()
    S->>M: Save to DB
    M-->>S: Success
    S->>K: Publish sync event
    S-->>A: Return response
    A-->>C: 201 Created
    
    Note over K,E: Background Sync
    K->>E: Consumer processes event
    E->>E: Index document
    
    Note over C,E: Search Flow
    C->>A: GET /starters/search?q=john
    A->>S: SearchStarters()
    S->>E: Query index
    E-->>S: Results
    S-->>A: Return data
    A-->>C: 200 OK
```

## Tech Stack

```mermaid
graph LR
    subgraph Backend
        Go[Go 1.21+]
        Gin[Gin Web Framework]
        GORM[GORM ORM]
    end
    
    subgraph Storage
        MySQL[MySQL 8.0]
        ES[Elasticsearch 8.x]
    end
    
    subgraph Messaging
        Kafka[Kafka 3.x]
        Sarama[Sarama Client]
    end
    
    subgraph Tools
        Swagger[Swagger/OpenAPI]
        Docker[Docker Compose]
    end
    
    Go --> Gin
    Go --> GORM
    Go --> Sarama
    GORM --> MySQL
    Sarama --> Kafka
    Kafka --> ES
```

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Ports available: `3000`, `3306`, `9092`, `9200`

### Run with Docker

```bash
# Start all services
make up

# Check status
make ps

# View logs
make logs-app

# Stop services
make down

# Clean everything (including volumes)
make clean
```

### Access Points

| Service | URL | Description |
|---------|-----|-------------|
| API | `http://localhost:3000/api/v1` | REST endpoints |
| Swagger | `http://localhost:3000/swagger/index.html` | API docs |
| Elasticsearch | `http://localhost:9200` | Search engine |
| Kibana | `http://localhost:5601` | ES visualization |

## Project Structure

```
.
├── cmd/api/                 # Application entry point (main.go)
├── internal/
│   ├── config/              # Environment configuration
│   ├── initialize/          # Bootstrap logic
│   ├── organization/        # Domain module
│   │   ├── domain/          # Entities & business rules
│   │   ├── application/     # Use cases & services
│   │   └── infrastructure/  # Repositories & adapters
│   └── shared/              # Shared infrastructure
│       ├── messagebroker/   # Kafka client
│       └── search/          # Elasticsearch
├── pkg/                     # Utilities (httpctx, response)
├── docs/                    # Swagger/OpenAPI specs
├── migrations/              # SQL migration scripts
├── docker-compose.yml       # Service definitions
└── Makefile                 # Build & deployment commands
```

## Key Features

```mermaid
mindmap
  root((Starter API))
    CRUD Operations
      Create starter
      Update info
      Delete records
      List with pagination
    Search
      Full-text search
      Field-specific filters
      Autocomplete support
      Fallback to MySQL
    Sync Mechanism
      Event-driven architecture
      Kafka message queue
      Background indexing
      Auto-retry on failure
    Organization
      Companies
      Business Units
      Departments
      Hierarchy management
```

## Configuration

### Environment Variables

```mermaid
graph LR
    subgraph Database
        DB_HOST
        DB_PORT
        DB_USER
        DB_PASS
        DB_NAME
    end
    
    subgraph Elasticsearch
        ES_ADDR[ES_ADDRESSES]
        ES_USER[ES_USERNAME]
        ES_PASS[ES_PASSWORD]
    end
    
    subgraph Kafka
        KAFKA_BROKERS
        KAFKA_TOPIC
        KAFKA_GROUP[CONSUMER_GROUP]
    end
    
    subgraph Server
        PORT[SERVER_PORT]
        LOG[LOG_LEVEL]
    end
    
    ENV[.env_dev] -.-> Database
    ENV -.-> Elasticsearch
    ENV -.-> Kafka
    ENV -.-> Server
```

**Note**: If Kafka or Elasticsearch is unavailable, the system gracefully degrades functionality.

## Event Flow

```mermaid
stateDiagram-v2
    [*] --> APIRequest
    APIRequest --> ValidateData
    ValidateData --> SaveMySQL
    SaveMySQL --> PublishEvent
    PublishEvent --> KafkaTopic
    
    KafkaTopic --> ConsumerPick
    ConsumerPick --> FetchLatest
    FetchLatest --> IndexES
    IndexES --> [*]
    
    ValidateData --> [*]: Invalid
    PublishEvent --> [*]: Kafka disabled
```

## Elasticsearch Index

```mermaid
graph TB
    subgraph Index Structure
        Index[starters index]
        Fields[Fields]
        Analyzers[Custom Analyzers]
    end
    
    Fields --> ID[id: long]
    Fields --> Domain[domain: text + keyword]
    Fields --> Email[email: text + keyword]
    Fields --> JobTitle[job_title: text]
    Fields --> FullText[full_text: searchable]
    
    Analyzers --> EdgeNgram[edge_ngram<br/>min: 2, max: 10]
    Analyzers --> Lowercase[lowercase]
    Analyzers --> ASCII[asciifolding]
    
    EdgeNgram -.->|Indexing| Domain
    Lowercase -.->|Both| Domain
    ASCII -.->|Both| Domain
    
    style Index fill:#005571,color:#fff
```

## API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/starters` | Create new starter |
| `GET` | `/api/v1/starters/:id` | Get starter details |
| `PUT` | `/api/v1/starters/:id` | Update starter |
| `DELETE` | `/api/v1/starters/:id` | Delete starter |
| `GET` | `/api/v1/starters` | List all starters |
| `GET` | `/api/v1/starters/search` | Search starters |
| `POST` | `/api/v1/starters/reindex` | Trigger full reindex |

## Development Workflow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Git as Git
    participant CI as CI/CD
    participant Docker as Docker
    participant App as Application

    Dev->>Git: Push code
    Git->>CI: Trigger build
    CI->>CI: Run tests
    CI->>Docker: Build image
    Docker->>App: Deploy container
    App->>App: Run migrations
    App->>App: Start services
    
    Note over App: Kafka Consumer<br/>ES Indexer<br/>API Server
```

## Monitoring & Operations

```mermaid
graph TD
    subgraph Health Checks
        API[API Health]
        MySQL[MySQL Connection]
        ES[ES Cluster Status]
        Kafka[Kafka Broker]
    end
    
    subgraph Logs
        AppLog[Application Logs]
        ESLog[ES Sync Logs]
        KafkaLog[Kafka Consumer Logs]
    end
    
    subgraph Metrics
        ReqCount[Request Count]
        IndexSpeed[Index Speed]
        QueueLag[Queue Lag]
    end
    
    API --> AppLog
    MySQL --> AppLog
    ES --> ESLog
    Kafka --> KafkaLog
    
    AppLog --> ReqCount
    ESLog --> IndexSpeed
    KafkaLog --> QueueLag
```

## Troubleshooting

| Issue | Check | Solution |
|-------|-------|----------|
| API not responding | `make logs-app` | Check port 3000 availability |
| Search not working | `curl localhost:9200` | Verify ES is running |
| Sync delays | `make logs-kafka` | Check Kafka consumer lag |
| Missing data | MySQL connection | Verify DB credentials in `.env_dev` |

## Performance Considerations

```mermaid
graph LR
    subgraph Optimization
        A[MySQL Indexes] --> B[Fast Queries]
        C[ES Shards] --> D[Parallel Search]
        E[Kafka Partitions] --> F[High Throughput]
        G[Connection Pooling] --> H[Reduced Latency]
    end
    
    B --> Performance[Better Performance]
    D --> Performance
    F --> Performance
    H --> Performance
    
    style Performance fill:#90EE90
```

## License

MIT

## Support

For issues or questions, please check the logs:
```bash
make logs-app     # Application logs
make logs-mysql   # Database logs
make logs-es      # Elasticsearch logs
make logs-kafka   # Kafka logs
```