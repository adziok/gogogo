package externalapi

type ExternalFeatureFlag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"value"`
}
