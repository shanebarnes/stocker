package main

import (
	"os"
	"testing"

	"github.com/shanebarnes/stocker/internal/stock/api"
	"github.com/stretchr/testify/assert"
)

func TestInitEnvVars(t *testing.T) {
	// Save current environment variables before modifying
	saveKey := api.GetApiKeyFromEnv()
	saveServer := api.GetApiServerFromEnv()

	// Restore original environment variable values after tests
	defer os.Setenv(api.ApiKeyEnvName, saveKey)
	defer os.Setenv(api.ApiServerEnvName, saveServer)

	os.Setenv(api.ApiKeyEnvName, "")
	os.Setenv(api.ApiServerEnvName, "")
	initEnvVars()
	assert.Equal(t, "", apiKey)
	assert.Equal(t, "", apiServer)

	os.Setenv(api.ApiKeyEnvName, "SomeApiKey")
	os.Setenv(api.ApiServerEnvName, "SomeApiServer")
	initEnvVars()
	assert.Equal(t, "SomeApiKey", apiKey)
	assert.Equal(t, "SomeApiServer", apiServer)
}