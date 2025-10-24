# Go Microservices - Starter & Notification Management

## Quick Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/kiin21/go-rest.git
   cd go-rest
   ```

2. **Configure environment variables**

   Both services include example configuration files:

   ```bash
   # Starter Service
   cp services/starter-service/.env_dev.example services/starter-service/.env_dev

   # Notification Service
   cp services/notification-service/.env_dev.example services/notification-service/.env_dev
   ```

3. **Start infrastructure services first**

   ```bash
   make infra-up
   # This starts: MySQL, MongoDB, Kafka, Zookeeper, Elasticsearch
   # Or: docker compose -f docker-compose.infra.yml up -d
   ```

   Wait for all services to be healthy. You can check with:

   ```bash
   docker compose -f docker-compose.infra.yml ps
   ```

4. **Start application services**

   ```bash
   make up
   # This starts: Starter Service, Notification Service
   # Or: docker compose up -d
   ```



5. **Access the application**
   - **Starter Service API**: `http://localhost:8080/api/v1`
   - **Swagger Documentation**: `http://localhost:8080/swagger/index.html`
   - **Notification Service API**: `http://localhost:8081/api/v1`
   - **Elasticsearch**: `http://localhost:9200`
   - **Kibana** (optional): `http://localhost:5601`

### Environment Configuration

**Starter Service** (`services/starter-service/.env_dev`):

- `DB_URI` - MySQL connection string
- `SERVER_PORT` - HTTP server port (default: 3000)
- `ELASTICSEARCH_ADDRESSES` - Elasticsearch URL
- `KAFKA_BROKERS` - Kafka broker addresses

**Notification Service** (`services/notification-service/.env_dev`):

- `MONGODB_URI` - MongoDB connection string
- `SERVER_PORT` - HTTP server port (default: 8081)
- `KAFKA_BROKERS` - Kafka broker addresses

## Testing

### Unit Tests

```bash
# With coverage
cd services/starter-service
go test ./internal/starter/domain/... ./internal/starter/application/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
### Integration Tests

### Testing

- `make test` – run unit tests
- `make test-starter` – test starter service only
- `make test-notification` – test notification service only
- `make test-integration` – run integration tests (requires Docker)
- `make test-all` – run unit + integration tests
