package providers

import (
	"certify/internal/acme/constants/provider_constants"
	"certify/internal/acme/providers/cloudflare"
	"certify/internal/acme/providers/websupport"
	"certify/internal/acme/zone_configuration"
	"certify/internal/configuration"
	"fmt"
	legoCertificate "github.com/go-acme/lego/v4/certificate"
)

type Provider interface {
	ObtainCertificate(configuration *configuration.Configuration, zoneConfiguration *zone_configuration.ZoneConfiguration) (*legoCertificate.Resource, error)
}

func GetProvider(provider provider_constants.Provider) Provider {
	switch provider {
	case provider_constants.ProviderCloudflare:
		return cloudflare.NewProvider()
	case provider_constants.ProviderWebsupport:
		return websupport.NewProvider()
	default:
		panic(fmt.Errorf("unknown provider: %s", provider))
	}
}
