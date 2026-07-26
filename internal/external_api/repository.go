package externalapi

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Operation[T any] struct {
	Data   T
	Tenant string
	User   string
}

type FeatureFlagExternalReposiotory interface {
	GetByTenantAndName(ctx context.Context, tenant string, name string) (ExternalFeatureFlag, error)
}

type FeatureFlagExternalPostgresReposiotory struct {
	db *pgxpool.Pool
}

func NewFeatureFlagExternalPostgresReposiotory(db *pgxpool.Pool) *FeatureFlagExternalPostgresReposiotory {
	return &FeatureFlagExternalPostgresReposiotory{
		db: db,
	}
}

func (r *FeatureFlagExternalPostgresReposiotory) GetByTenantAndName(ctx context.Context, tenant string, name string) (ExternalFeatureFlag, error) {
	query := `
	SELECT name, enabled
	FROM public.feature_flags
	WHERE tenant = $1
	AND name = $2;
	`
	slog.Debug("args:", name, tenant)

	var f ExternalFeatureFlag
	err := r.db.QueryRow(ctx, query, tenant, name).Scan(
		&f.Name,
		&f.Enabled)

	if err != nil {
		return ExternalFeatureFlag{}, fmt.Errorf("postgres failed to find feature flag: %w", err)
	}

	return f, nil
}
