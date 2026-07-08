package feature_flags

import "context"

type Operation[T any] struct {
	Data   T
	Tenant string
	User   string
}

type FeatureFlagRepository interface {
	Create(ctx context.Context, f Operation[CreateFeatureFlag]) error
	Update(ctx context.Context, f Operation[UpdateFeatureFlag]) error
	DeleteById(ctx context.Context, f Operation[DeleteFeatureFlag]) error
	GetByTenant(ctx context.Context, tenant string) ([]FeatureFlag, error)
}
