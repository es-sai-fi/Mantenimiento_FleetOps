package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a connection pool to the PostgreSQL database.
// It uses pgx v5's pgxpool for connection pooling and management.
//
// SAD Reference: ADR-2 — "PostgreSQL mediante Supabase"
// The connection string points to a Supabase-hosted PostgreSQL instance,
// but we use the standard pgx driver (not Supabase SDK) per confirmed decision FC-2.
func NewPostgresPool(ctx context.Context, databaseURL string, maxConns int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
