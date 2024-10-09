package acme

import (
	"certify/internal/acme/providers"
	"certify/internal/acme/zone_configuration"
	"certify/internal/certificates"
	"certify/internal/configuration"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func HandleZone(config *configuration.Configuration, zoneConfiguration *zone_configuration.ZoneConfiguration) error {
	log.Printf("Handling zone: %s", zoneConfiguration.UniqueIdentifier)

	certificateDirectoryPath := path.Join(config.CertificatesPath, zoneConfiguration.UniqueIdentifier)
	if err := os.MkdirAll(certificateDirectoryPath, 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}

	certificateExpirationDays, err := certificates.GetExpirationDays(certificateDirectoryPath)
	log.Printf("Certificate expiration is %d days from now", certificateExpirationDays)
	if err != nil {
		return fmt.Errorf("failed to get certificate expiration days: %w", err)
	}

	if certificateExpirationDays > zoneConfiguration.RenewalDays {
		log.Printf("Certificate expiration is more than %d days from now, skipping", zoneConfiguration.RenewalDays)
		return nil
	}

	log.Printf("Certificate expiration is less than %d days from now, renewing", zoneConfiguration.RenewalDays)

	provider := providers.GetProvider(zoneConfiguration.Provider)
	certificate, err := provider.ObtainCertificate(config, zoneConfiguration)
	if err != nil {
		return fmt.Errorf("failed to obtain certificate: %w", err)
	}

	if err = certificates.SaveCertificate(certificateDirectoryPath, certificate, zoneConfiguration); err != nil {
		return fmt.Errorf("failed to save certificate: %w", err)
	}

	return nil
}

func GetZones(path string) []*zone_configuration.ZoneConfiguration {
	zoneConfigurationPaths, err := getZoneConfigurationsInDirectory(path)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get zone configuration files: %w", err))
	}

	var zoneConfigurations []*zone_configuration.ZoneConfiguration
	for _, zoneConfigurationPath := range zoneConfigurationPaths {
		zoneConfiguration, err := zone_configuration.ReadZoneConfiguration(zoneConfigurationPath)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to read zone configuration: %w", err))
		}

		zoneConfigurations = append(zoneConfigurations, zoneConfiguration)
	}

	if err := areZoneConfigurationsIdentifiersUnique(zoneConfigurations); err != nil {
		log.Fatal(err)
	}

	return zoneConfigurations
}

func areZoneConfigurationsIdentifiersUnique(zoneConfigurations []*zone_configuration.ZoneConfiguration) error {
	uniqueIdentifiers := make(map[string]bool)
	for _, zoneConfiguration := range zoneConfigurations {
		if _, ok := uniqueIdentifiers[zoneConfiguration.UniqueIdentifier]; ok {
			return fmt.Errorf("there are multiple zone configurations with the same unique identifier: %s", zoneConfiguration.UniqueIdentifier)
		}

		uniqueIdentifiers[zoneConfiguration.UniqueIdentifier] = true
	}

	return nil
}

func getZoneConfigurationsInDirectory(zoneDirectoryPath string) ([]string, error) {
	stat, err := os.Stat(zoneDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones directory: %w", err)
	}

	if stat.IsDir() == false {
		return nil, fmt.Errorf("zones path is not a directory: %s", zoneDirectoryPath)
	}

	files, err := os.ReadDir(zoneDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read zones directory: %w", err)
	}

	var zoneConfigurationPaths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Ignore files that are not .yaml or .yml
		extension := filepath.Ext(file.Name())
		if extension != ".yaml" && extension != ".yml" {
			continue
		}

		zoneConfigurationPaths = append(zoneConfigurationPaths, path.Join(zoneDirectoryPath, file.Name()))
	}

	return zoneConfigurationPaths, nil
}
