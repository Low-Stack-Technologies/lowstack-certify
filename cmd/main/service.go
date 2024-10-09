package main

import (
	"certify/internal/acme"
	"certify/internal/configuration"
	"log"
	"time"
)

func main() {
	log.Println("Low-Stack Certify is running!")

	// Load the configuration
	config := configuration.GetConfiguration()

	for {
		// Re-load the configuration
		config = configuration.GetConfiguration()

		// Get all zones and handle them
		zones := acme.GetZones(config.ZonesPath)
		for _, zone := range zones {
			if err := acme.HandleZone(config, zone); err != nil {
				log.Printf("Failed to handle zone: %s\n%s", zone.UniqueIdentifier, err)
			}
		}

		// Break if the application should not run periodically
		if !config.RuntimeConfiguration.RunPeriodically {
			break
		}

		// Sleep for the specified amount of time
		time.Sleep(time.Duration(config.RuntimeConfiguration.PeriodMinutes) * time.Minute)
	}

	log.Println("Low-Stack Certify is exiting!")
}
