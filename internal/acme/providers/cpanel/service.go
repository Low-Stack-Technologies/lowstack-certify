package cpanel

import (
	"certify/internal/acme/providers/provider_utils"
	"certify/internal/acme/zone_configuration"
	"certify/internal/configuration"
	"fmt"
	legoCertificate "github.com/go-acme/lego/v4/certificate"
	cpanelChallenge "github.com/go-acme/lego/v4/providers/dns/cpanel"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (p Provider) ObtainCertificate(configuration *configuration.Configuration, zoneConfiguration *zone_configuration.ZoneConfiguration) (*legoCertificate.Resource, error) {
	acmeUser, exists, err := provider_utils.GetACMEUser(configuration, zoneConfiguration.IdentityEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	acmeClient, err := provider_utils.GetACMEClient(acmeUser, configuration, zoneConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to get ACME client: %w", err)
	}

	defer provider_utils.UnsetEnvironmentVariable("CPANEL_USERNAME")
	if err = provider_utils.SetEnvironmentVariable("CPANEL_USERNAME", zoneConfiguration, "username"); err != nil {
		return nil, fmt.Errorf("failed to set CPanel username: %w", err)
	}

	defer provider_utils.UnsetEnvironmentVariable("CPANEL_TOKEN")
	if err = provider_utils.SetEnvironmentVariable("CPANEL_TOKEN", zoneConfiguration, "token"); err != nil {
		return nil, fmt.Errorf("failed to set CPanel API token: %w", err)
	}

	defer provider_utils.UnsetEnvironmentVariable("CPANEL_BASE_URL")
	if err = provider_utils.SetEnvironmentVariable("CPANEL_BASE_URL", zoneConfiguration, "base_url"); err != nil {
		return nil, fmt.Errorf("failed to set CPanel base URL: %w", err)
	}

	// Optional set CPanel mode
	if zoneConfiguration.ProviderOptions["mode"] != "" {
		defer provider_utils.UnsetEnvironmentVariable("CPANEL_MODE")
		if err = provider_utils.SetEnvironmentVariable("CPANEL_MODE", zoneConfiguration, "mode"); err != nil {
			return nil, fmt.Errorf("failed to set CPanel mode: %w", err)
		}
	}

	dnsProvider, err := cpanelChallenge.NewDNSProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS provider: %w", err)
	}

	if err := acmeClient.Challenge.SetDNS01Provider(dnsProvider); err != nil {
		return nil, fmt.Errorf("failed to set DNS provider: %w", err)
	}

	// If the user does not exist, register them
	if !exists {
		if err := provider_utils.RegisterACMEUser(acmeClient, acmeUser); err != nil {
			return nil, fmt.Errorf("failed to register user: %w", err)
		}
	}

	return provider_utils.ObtainACMECertificate(acmeClient, zoneConfiguration)
}
