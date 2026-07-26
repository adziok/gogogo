package auth

import (
	"fmt"
	"os"
)

type AuthConfig struct {
	Domain      string
	Audience    string
	AudienceApi string
}

func LoadAuthConfig() (*AuthConfig, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN environment variable required")
	}

	audience := os.Getenv("AUTH0_AUDIENCE")
	if audience == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE environment variable required")
	}

	audienceApi := os.Getenv("AUTH0_AUDIENCE_API")
	if audience == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE environment variable required")
	}

	return &AuthConfig{
		Domain:      domain,
		Audience:    audience,
		AudienceApi: audienceApi,
	}, nil
}
