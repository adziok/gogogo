package feature_flags

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (r *FeatureFlagPostgresRepository) Create(ctx context.Context, f Operation[CreateFeatureFlag]) error {
	query := `
	INSERT INTO public.feature_flags (description, "name", enabled, tenant, created_at, created_by, updated_at, updated_by, id)
	VALUES ($1, $2, $3, $4, timezone('utc'::text, now()), $5, timezone('utc'::text, now()), $6, $7);
	`
	_, err := r.db.Exec(ctx, query, f.Data.Description, f.Data.Name, f.Data.Enabled, f.Tenant, f.User, f.User, uuid.New().String())

	if err != nil {
		return fmt.Errorf("postgres failed to save feature flag: %w", err)
	}

	return nil
}

// UPDATE public.feature_flags SET description=”, "name"=”, enabled=false, tenant=”, created_at=timezone('utc'::text, now()), created_by=”, updated_at=timezone('utc'::text, now()), updated_by=” WHERE id=?;
func (r *FeatureFlagPostgresRepository) Update(ctx context.Context, f Operation[UpdateFeatureFlag]) error {
	query := `
	UPDATE public.feature_flags
	SET description = $1, enabled = $2, updated_by = $3, updated_at = timezone('utc'::text, now())
	WHERE id = $4 AND tenant = $5;
	`
	_, err := r.db.Exec(ctx, query, f.Data.Description, f.Data.Enabled, f.User, f.Data.ID, f.Tenant)

	if err != nil {
		return fmt.Errorf("postgres failed to update feature flag: %w", err)
	}

	return nil
}

func (r *FeatureFlagPostgresRepository) DeleteById(ctx context.Context, f Operation[DeleteFeatureFlag]) error {
	query := `
	DELETE FROM public.feature_flags
	WHERE id = $1 AND tenant = $2;
	`
	_, err := r.db.Exec(ctx, query, f.Data.ID, f.Tenant)

	if err != nil {
		return fmt.Errorf("postgres failed to delete feature flag: %w", err)
	}

	return nil
}

func (r *FeatureFlagPostgresRepository) GetByTenant(ctx context.Context, tenant string) ([]FeatureFlag, error) {
	query := `
	SELECT id, name, description, tenant, enabled, created_at, created_by, updated_at, updated_by
	FROM public.feature_flags
	WHERE tenant = $1
	ORDER BY created_at DESC;
	`
	result, err := r.db.Query(ctx, query, tenant)

	if err != nil {
		return nil, fmt.Errorf("postgres failed to delete feature flag: %w", err)
	}

	defer result.Close()

	flags := make([]FeatureFlag, 0)

	for result.Next() {
		var f FeatureFlag

		// Mapujemy kolumny z bieżącego wiersza na strukturę
		if err := result.Scan(
			&f.ID,
			&f.Name,
			&f.Description,
			&f.Tenant,
			&f.Enabled,
			&f.CreatedAt,
			&f.CreatedBy,
			&f.UpdatedAt,
			&f.UpdatedBy,
		); err != nil {
			return nil, fmt.Errorf("postgres failed to scan feature flag row: %w", err)
		}

		// Dodajemy sparsowany obiekt do naszej listy
		flags = append(flags, f)
	}
	// 5. Po zakończeniu pętli zawsze sprawdzamy, czy w trakcie iteracji nie wystąpił błąd bazy
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("postgres rows error: %w", err)
	}

	return flags, nil
}
