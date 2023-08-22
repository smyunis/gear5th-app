package mongotestdoubles

import (
	"bytes"
	_ "embed"

	"github.com/joho/godotenv"
)

type TestEnvConfigurationProvider struct {
	env map[string]string
}

//go:embed .env.test
var envFile []byte

func NewTestEnvConfigurationProvider() TestEnvConfigurationProvider {

	env, err := godotenv.Parse(bytes.NewReader(envFile))
	if err != nil {
		panic(err)
	}
	return TestEnvConfigurationProvider{
		env,
	}
}

func (e TestEnvConfigurationProvider) Get(key string, defaultValue string) string {
	envVar := e.env[key]
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
