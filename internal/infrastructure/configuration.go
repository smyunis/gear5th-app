package infrastructure

import (
	"os"
	// "github.com/joho/godotenv"
)

type ConfigurationProvider interface {
	Get(key, defaultValue string) string
}

// func init() {
// 	err := godotenv.Load("config/.env.dev", "config/.env.prod")
// 	if err != nil {
// 		panic("could not load config file ./config/.env.*")
// 	}
// }

type EnvConfigurationProvider struct{}

func NewEnvConfigurationProvider() EnvConfigurationProvider {
	return EnvConfigurationProvider{}
}

func (e EnvConfigurationProvider) Get(key string, defaultValue string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
