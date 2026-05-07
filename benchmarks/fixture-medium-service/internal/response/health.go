package response

type Health struct {
	Status      string         `json:"status"`
	Service     string         `json:"service"`
	Environment string         `json:"environment"`
	Region      string         `json:"region"`
	BuildNumber int            `json:"build_number"`
	Dependency  map[string]any `json:"dependency"`
	Timestamp   string         `json:"timestamp"`
}
