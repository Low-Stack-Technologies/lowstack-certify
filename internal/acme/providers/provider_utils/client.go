package provider_utils

import (
	"certify/internal/acme/zone_configuration"
	"certify/internal/configuration"
	"fmt"
	legoCertificate "github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
)

func GetACMEClient(acmeUser *User, configuration *configuration.Configuration, zoneConfiguration *zone_configuration.ZoneConfiguration) (*lego.Client, error) {
	// Create a new ACME config
	acmeConfig := lego.NewConfig(acmeUser)
	acmeConfig.CADirURL = configuration.CAURL
	acmeConfig.Certificate.KeyType = zoneConfiguration.KeyType

	// Create ACME client
	acmeClient, err := lego.NewClient(acmeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACME client: %w", err)
	}

	return acmeClient, nil
}

func ObtainACMECertificate(acmeClient *lego.Client, zoneConfiguration *zone_configuration.ZoneConfiguration) (*legoCertificate.Resource, error) {
	request := legoCertificate.ObtainRequest{
		Domains: zoneConfiguration.Hostnames,
		Bundle:  true,
	}
	certificate, err := acmeClient.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain certificate: %w", err)
	}

	return certificate, nil
}
