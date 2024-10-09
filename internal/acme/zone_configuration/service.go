package zone_configuration

import (
	"certify/internal/acme/constants/provider_constants"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"gopkg.in/yaml.v3"
	"os"
)

type ZoneConfiguration struct {
	UniqueIdentifier string                      `yaml:"unique_identifier"`
	Hostnames        []string                    `yaml:"hostnames"`
	IdentityEmail    string                      `yaml:"identity_email"`
	RenewalDays      int                         `yaml:"renewal_days"`
	Provider         provider_constants.Provider `yaml:"provider"`
	ProviderOptions  map[string]string           `yaml:"provider_options"`
	KeyType          certcrypto.KeyType          `yaml:"key_type"`
}

func ReadZoneConfiguration(path string) (*ZoneConfiguration, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read zone configuration @ %s: %w", path, err)
	}

	var zone ZoneConfiguration
	err = yaml.Unmarshal(contents, &zone)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal zone configuration: %w", err)
	}

	return &zone, nil
}
