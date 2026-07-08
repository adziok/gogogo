package feature_flags

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type FeatureFlag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tenant      string    `json:"tenant"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

type CreateFeatureFlag struct {
	Name        string `json:"name" validate:"required,alphanum,min=3"`
	Description string `json:"description" validate:"required,max=100"`
	Enabled     bool   `json:"enabled" validate:"boolean"`
}

type DeleteFeatureFlag struct {
	ID string `json:"id" validate:"required,uuid"`
}

type UpdateFeatureFlag struct {
	ID          string `json:"id" validate:"required,uuid"`
	Description string `json:"description" validate:"required,max=100"`
	Enabled     bool   `json:"enabled" validate:"boolen"`
}

var validate = validator.New()
