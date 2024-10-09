package configuration

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	ZonesPath        string `yaml:"zones_path"`
	CertificatesPath string `yaml:"certificates_path"`
	CAURL            string `yaml:"ca_url"`
}

const defaultConfiguration = `
# This is the default configuration file for the Low-Stack Certify application.
# You can change all values here to customize the application to your needs.

# The path to the directory containing the zone configuration files.
# This can be a relative or absolute path.
zones_path: "/zones"

# The path to the directory where certificates will be stored.
# This can be a relative or absolute path.
certificates_path: "/certificates"

# This is the URL of the CA directory that will be used to sign certificates.
#ca_url: "https://acme-staging-v02.api.letsencrypt.org/directory" # Use this for testing
ca_url: "https://acme-v02.api.letsencrypt.org/directory"
`

func GetConfiguration() *Configuration {
	configPath := GetConfigurationPath()
	if err := WriteDefaultConfigurationIfNotExists(configPath); err != nil {
		log.Fatal(fmt.Errorf("failed to write default configuration: %w", err))
	}

	config, err := ReadConfiguration(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read configuration: %w", err))
	}

	return config
}

func ReadConfiguration(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Configuration
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteDefaultConfigurationIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) == false {
		return nil
	}

	// This will fail if the default configuration
	// does not match the expected structure
	ValidateDefaultConfiguration()

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}
	defer file.Close()

	trimmedDefaultConfiguration := strings.TrimSpace(defaultConfiguration)
	if _, err := file.WriteString(trimmedDefaultConfiguration); err != nil {
		return fmt.Errorf("failed to write default configuration: %w", err)
	}

	return nil
}

func GetConfigurationPath() string {
	if path, ok := os.LookupEnv("CUSTOM_CONFIGURATION_PATH"); ok {
		return path
	}

	return "config/config.yaml"
}

func ValidateDefaultConfiguration() {
	if err := yaml.Unmarshal([]byte(defaultConfiguration), &Configuration{}); err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal default configuration: %w", err))
	}
}
