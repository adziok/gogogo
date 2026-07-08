package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// NewClickHouseClient inicjalizuje połączenie z ClickHouse przez natywny protokół TCP
func NewClickHouseClient(ctx context.Context, connStr string) (clickhouse.Conn, error) {
	// 1. Parsujemy connection string
	options, err := clickhouse.ParseDSN(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse clickhouse options: %w", err)
	}

	// 2. Otwieramy połączenie
	conn, err := clickhouse.Open(options)
	if err != nil {
		return nil, fmt.Errorf("unable to open clickhouse connection: %w", err)
	}

	// 3. Sprawdzamy czy baza odpowiada
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse ping failed: %w", err)
	}

	slog.Info("Successfully connected to ClickHouse")
	return conn, nil
}
