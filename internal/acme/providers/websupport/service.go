package websupport

import (
	"certify/internal/acme/providers/provider_utils"
	"certify/internal/acme/zone_configuration"
	"certify/internal/configuration"
	"fmt"
	legoCertificate "github.com/go-acme/lego/v4/certificate"
	websupportChallenge "github.com/go-acme/lego/v4/providers/dns/websupport"
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

	defer provider_utils.UnsetEnvironmentVariable("WEBSUPPORT_API_KEY")
	if err = provider_utils.SetEnvironmentVariable("WEBSUPPORT_API_KEY", zoneConfiguration, "api_key"); err != nil {
		return nil, fmt.Errorf("failed to set Websupport API key: %w", err)
	}

	defer provider_utils.UnsetEnvironmentVariable("WEBSUPPORT_SECRET")
	if err = provider_utils.SetEnvironmentVariable("WEBSUPPORT_SECRET", zoneConfiguration, "api_secret"); err != nil {
		return nil, fmt.Errorf("failed to set Websupport API secret: %w", err)
	}

	dnsProvider, err := websupportChallenge.NewDNSProvider()
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
