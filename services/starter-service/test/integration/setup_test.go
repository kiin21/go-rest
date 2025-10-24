package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/initialize"
	initStarter "github.com/kiin21/go-rest/services/starter-service/internal/initialize/starter"
	persistentMySQL "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/persistence/repository/mysql"
)

// TestEnv holds all test infrastructure
type TestEnv struct {
	Router         *gin.Engine
	DB             *gorm.DB
	MySQLContainer testcontainers.Container
	MySQLConnStr   string
	Cleanup        func()
}

// SetupTestEnvironment initializes all test infrastructure
func SetupTestEnvironment(t *testing.T) *TestEnv {
	ctx := context.Background()

	// Start MySQL container with testcontainers
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("intern_app_test"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server").
				WithOccurrence(1).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %v", err)
	}

	// Get connection string
	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get MySQL host: %v", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("Failed to get MySQL port: %v", err)
	}

	connStr := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/intern_app_test?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())

	// Wait a bit for MySQL to be fully ready
	time.Sleep(2 * time.Second)

	// Connect to database
	db, err := gorm.Open(mysqlDriver.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db, connStr); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize test router and dependencies
	router := setupRouter(db)

	env := &TestEnv{
		Router:         router,
		DB:             db,
		MySQLContainer: mysqlContainer,
		MySQLConnStr:   connStr,
		Cleanup: func() {
			if err := mysqlContainer.Terminate(ctx); err != nil {
				log.Printf("Failed to terminate MySQL container: %v", err)
			}
		},
	}

	return env
}

// setupRouter creates a Gin router with all dependencies for testing
func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Create mock/no-op producers for testing
	notificationProducer := &NoOpProducer{}
	syncProducer := &NoOpSyncProducer{}

	// Setup repositories
	requestURLResolver := httputil.NewRequestURLResolver()
	starterRepo := persistentMySQL.NewStarterRepository(db)
	businessUnitRepo := persistentMySQL.NewBusinessUnitRepository(db)
	departmentRepo := persistentMySQL.NewDepartmentRepository(db)

	// Initialize handlers
	orgHandler := initStarter.InitOrganization(
		starterRepo,
		departmentRepo,
		businessUnitRepo,
		notificationProducer,
	)

	starterHandler, _, _ := initStarter.InitStarter(
		starterRepo,
		departmentRepo,
		businessUnitRepo,
		nil, // No Elasticsearch for basic tests
		syncProducer,
	)

	// Initialize router
	router := initialize.InitRouter(
		"debug",
		requestURLResolver,
		orgHandler,
		starterHandler,
	)

	return router
}

// runMigrations executes all SQL migration files
func runMigrations(db *gorm.DB, connStr string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Get the migrations directory
	migrationsPath := "../../migrations"
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Execute each migration file in order
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		migrationPath := filepath.Join(migrationsPath, file.Name())
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Execute the migration
		if err := executeSQLFile(sqlDB, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}

		log.Printf("Executed migration: %s", file.Name())
	}

	return nil
}

// executeSQLFile executes SQL statements from a file content
// It splits the content into individual statements and executes them one by one
func executeSQLFile(db *sql.DB, content string) error {
	// Split SQL file into statements (naive split by semicolon)
	// This is a simple approach that works for most cases
	// For production, consider using a proper SQL parser
	statements := splitSQLStatements(content)

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement %d: %w\nStatement: %s", i+1, err, stmt[:min(len(stmt), 200)])
		}
	}

	return nil
}

// splitSQLStatements splits SQL content into individual statements
// It handles DELIMITER changes for stored procedures
func splitSQLStatements(content string) []string {
	var statements []string
	
	// Remove DELIMITER commands and handle stored procedures
	content = preprocessSQL(content)
	
	var current strings.Builder
	var inString bool
	var stringChar rune
	var inProcedure bool

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip comments and empty lines
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}

		// Detect stored procedure boundaries
		upperLine := strings.ToUpper(trimmed)
		if strings.Contains(upperLine, "CREATE PROCEDURE") || strings.Contains(upperLine, "CREATE FUNCTION") {
			inProcedure = true
		}

		for i, ch := range line {
			if ch == '\'' || ch == '"' || ch == '`' {
				if !inString {
					inString = true
					stringChar = ch
				} else if ch == stringChar {
					// Check if it's escaped
					if i > 0 && line[i-1] != '\\' {
						inString = false
					}
				}
			}

			current.WriteRune(ch)

			// For procedures, look for END$$ or similar patterns
			if inProcedure {
				// Check if we've reached the end of the procedure
				currentStr := strings.TrimSpace(current.String())
				if strings.HasSuffix(strings.ToUpper(currentStr), "END$$") || 
				   (strings.HasSuffix(strings.ToUpper(currentStr), "END") && strings.Contains(currentStr, "$$")) {
					// Replace $$ with ; for execution
					stmt := strings.ReplaceAll(currentStr, "$$", "")
					stmt = strings.TrimSpace(stmt)
					if stmt != "" {
						statements = append(statements, stmt)
					}
					current.Reset()
					inProcedure = false
					continue
				}
			}

			// Check for statement terminator (semicolon outside of strings and not in procedures)
			if ch == ';' && !inString && !inProcedure {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" && stmt != ";" {
					statements = append(statements, stmt)
				}
				current.Reset()
			}
		}

		// Add newline to preserve formatting
		if current.Len() > 0 {
			current.WriteRune('\n')
		}
	}

	// Add the last statement if there's any content left
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" && stmt != ";" {
			// Handle procedure endings
			stmt = strings.ReplaceAll(stmt, "$$", "")
			stmt = strings.TrimSpace(stmt)
			if stmt != "" {
				statements = append(statements, stmt)
			}
		}
	}

	return statements
}

// preprocessSQL removes DELIMITER commands and cleans up the SQL
func preprocessSQL(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)
		
		// Skip DELIMITER commands
		if strings.HasPrefix(upper, "DELIMITER") {
			continue
		}
		
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CleanupDatabase removes all data from tables (for test isolation)
// It keeps the base seed data from migrations (companies, business units, base departments/starters)
func CleanupDatabase(t *testing.T, db *gorm.DB) {
	// Disable foreign key checks temporarily
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Clean only test-created data, keeping the base seed data
	// Base seed data:
	// - Companies: id 1
	// - Business Units: ids 1-4
	// - Departments: ids 1-86 (from migrations)
	// - Starters: ids 1-13 (from migrations)
	
	// Delete only starters created during tests (id > 13)
	if err := db.Exec("DELETE FROM starters WHERE id > 13").Error; err != nil {
		t.Logf("Warning: failed to clean starters: %v", err)
	}

	// Delete only departments created during tests (id > 86)
	// Departments up to id 86 are from migrations
	if err := db.Exec("DELETE FROM departments WHERE id > 86").Error; err != nil {
		t.Logf("Warning: failed to clean departments: %v", err)
	}

	// Don't delete business_units or companies - they are needed by tests
	// Tests reference business_unit_id = 1, 2, 3, 4 from migrations

	// Re-enable foreign key checks
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// NoOpProducer is a mock Kafka producer that does nothing (for testing without Kafka)
type NoOpProducer struct{}

func (p *NoOpProducer) SendNotification(event *events.Event) error {
	// No-op: do nothing for tests
	return nil
}

func (p *NoOpProducer) Close() error {
	return nil
}

// NoOpSyncProducer is a mock sync producer
type NoOpSyncProducer struct{}

func (p *NoOpSyncProducer) SendSyncEvent(event *events.Event) error {
	// No-op: do nothing for tests
	return nil
}

func (p *NoOpSyncProducer) Close() error {
	return nil
}

