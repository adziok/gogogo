package feature_flags

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeatureFlagPostgresRepository struct {
	db *pgxpool.Pool
}

func NewFeatureFlagPostgresRepository(db *pgxpool.Pool) *FeatureFlagPostgresRepository {
	return &FeatureFlagPostgresRepository{
		db: db,
	}
}

func (r *FeatureFlagPostgresRepository) Create(ctx context.Context, f CreateFeatureFlag, tenant string, creator string) error {
	query := `
	INSERT INTO public.feature_flags (description, "name", enabled, tenant, created_at, created_by, updated_at, updated_by, id)
	VALUES ($1, $2, $3, $4, timezone('utc'::text, now()), $5, timezone('utc'::text, now()), $6, $7);
	`
	_, err := r.db.Exec(ctx, query, f.Description, f.Name, f.Enabled, tenant, creator, creator, uuid.New().String())

	if err != nil {
		return fmt.Errorf("postgres failed to save feature flag: %w", err)
	}

	return nil
}

// TODO add tenant check
func (r *FeatureFlagPostgresRepository) GetByID(ctx context.Context, id string) (*FeatureFlag, error) {
	query := `
		SELECT *
		FROM feature_flags
		WHERE id = $1;
	`

	var f FeatureFlag

	err := r.db.QueryRow(ctx, query, id).Scan(&f.ID, &f.Enabled, &f.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("postgres failed to get feature flag by id: %w", err)
	}

	return &f, nil
}

func (r *FeatureFlagPostgresRepository) GetByTenant(ctx context.Context, tenant string) (*[]FeatureFlag, error) {
	panic("unimplemented")
}

func (r *FeatureFlagPostgresRepository) Update(ctx context.Context, f FeatureFlag) error {
	panic("unimplemented")
}
