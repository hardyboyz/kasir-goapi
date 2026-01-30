package config

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(connectionString string) (*pgxpool.Pool, error) {
	// Disable prepared statement caching by adding parameter to connection string
	if !strings.Contains(connectionString, "prepared_statement_cache_mode") {
		if strings.Contains(connectionString, "?") {
			connectionString += "&prepared_statement_cache_mode=disable"
		} else {
			connectionString += "?prepared_statement_cache_mode=disable"
		}
	}

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// Disable prepared statement caching to avoid conflicts
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Set connection pool settings
	config.MaxConns = 25
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	DB = pool
	return pool, nil
}
