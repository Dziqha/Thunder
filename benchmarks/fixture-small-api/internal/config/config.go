package config

import "os"

type Config struct {
	ServiceName string
	Environment string
}

func Load() Config {
	name := os.Getenv("SERVICE_NAME")
	if name == "" {
		name = "fixture-small-api"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	return Config{
		ServiceName: name,
		Environment: env,
	}
}
