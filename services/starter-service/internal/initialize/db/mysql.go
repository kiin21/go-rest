package db

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kiin21/go-rest/services/starter-service/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn, err := buildDSN(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	// GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var gormErr error
	DB, gormErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if gormErr != nil {
		return nil, fmt.Errorf("failed to connect database: %w", gormErr)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// Connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully.")
	return DB, nil
}

func buildDSN(uri string) (string, error) {
	if uri == "" {
		return "", fmt.Errorf("database URI is not configured")
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("invalid database URI: %w", err)
	}

	if parsed.Scheme != "mysql" {
		return "", fmt.Errorf("unsupported database scheme: %s", parsed.Scheme)
	}

	username := parsed.User.Username()
	if username == "" {
		return "", fmt.Errorf("database URI must include username")
	}

	password, _ := parsed.User.Password()
	host := parsed.Host
	if host == "" {
		return "", fmt.Errorf("database URI must include host")
	}

	database := strings.TrimPrefix(parsed.Path, "/")
	if database == "" {
		return "", fmt.Errorf("database URI must include database name")
	}

	query := parsed.RawQuery
	if query == "" {
		query = "charset=utf8mb4&parseTime=True&loc=Local"
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", username, password, host, database, query), nil
}
