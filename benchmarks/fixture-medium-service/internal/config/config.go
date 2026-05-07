package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName string
	Environment string
	Region      string
	BuildNumber int
}

func Load() Config {
	name := getenv("SERVICE_NAME", "fixture-medium-service")
	env := getenv("APP_ENV", "dev")
	region := getenv("APP_REGION", "local")
	buildNumber, _ := strconv.Atoi(getenv("BUILD_NUMBER", "1"))

	return Config{
		ServiceName: strings.TrimSpace(name),
		Environment: strings.TrimSpace(env),
		Region:      strings.TrimSpace(region),
		BuildNumber: buildNumber,
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
