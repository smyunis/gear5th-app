package infrastructure

import "os"

type ConfigurationProvider interface {
	Get(key, defaultValue string) string
}

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
