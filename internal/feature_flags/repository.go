package feature_flags

import "context"

type FeatureFlagRepository interface {
	Create(ctx context.Context, f CreateFeatureFlag, tenant string, creator string) error
	Update(ctx context.Context, f FeatureFlag) error
	GetByID(ctx context.Context, id string) (*FeatureFlag, error)
	GetByTenant(ctx context.Context, tenant string) (*[]FeatureFlag, error)
}
