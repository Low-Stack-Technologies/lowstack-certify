package provider_utils

import (
	"certify/internal/acme/zone_configuration"
	"fmt"
	"os"
)

func SetEnvironmentVariable(key string, zoneConfiguration *zone_configuration.ZoneConfiguration, optionKey string) error {
	value, ok := zoneConfiguration.ProviderOptions[optionKey]
	if !ok {
		return fmt.Errorf("no %s provided in configuration", optionKey)
	}

	err := os.Setenv(key, value)
	if err != nil {
		return fmt.Errorf("failed to set environment variable (%s): %w", key, err)
	}

	return nil
}

func UnsetEnvironmentVariable(key string) {
	err := os.Unsetenv(key)
	if err != nil {
		panic(fmt.Errorf("failed to unset environment variable (%s): %w", key, err))
	}
}
