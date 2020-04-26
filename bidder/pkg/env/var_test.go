package env_test

import (
	"auctioneer/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMandatory(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"NAME":          "service",
		"TIMEOUT_IN_MS": "1000",
	})
	defer cleanup()

	vars := &env.Vars{}

	assert.Equal(t, "service", vars.Mandatory("NAME"), "Mandatory string not read from environment variable")
	assert.Equal(t, 1000, vars.MandatoryInt("TIMEOUT_IN_MS"), "Mandatory integer not read from environment variable")
}

func TestMandatoryWhenMissing(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"PUSH_TEMPLATE": "{}",
	})
	defer cleanup()

	vars := &env.Vars{}

	vars.Mandatory("TIMEOUT_IN_MS")
	vars.Mandatory("REMOTE_ENDPOINT")

	assert.EqualError(t, vars.Error(), "missing mandatory configuration: TIMEOUT_IN_MS, REMOTE_ENDPOINT")
}

func TestMandatoryWhenMalformed(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"TIMEOUT_IN_MS": "1 minute",
	})
	defer cleanup()

	vars := &env.Vars{}
	vars.MandatoryInt("TIMEOUT_IN_MS")

	assert.EqualError(t, vars.Error(), `malformed configuration: mandatory TIMEOUT_IN_MS (value="1 minute") is not a number`)
}

func TestOptional(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"NAME":        "service",
		"LOG_LEVEL":   "error",
		"APP_PORT":    "80",
		"SMS_ENABLED": "true",
	})
	defer cleanup()

	vars := &env.Vars{}
	assert.Equal(t, "service", vars.Optional("NAME", "fallback"), "Optional string not read from environment variable")
	assert.Equal(t, 80, vars.OptionalInt("APP_PORT", 8080), "Optional integer not read from environment variable")

	require.NoError(t, vars.Error(), "Unexpected error")
}

func setEnvVars(t *testing.T, vars map[string]string) func() {
	for k, v := range vars {
		err := os.Setenv(k, v)
		require.NoError(t, err, "Could not set environment variable with name=%q and value=%q", k, v)
	}

	return func() {
		for k := range vars {
			err := os.Unsetenv(k)
			require.NoError(t, err, "Could not unset environment variable with name=%q", k)
		}
	}
}
