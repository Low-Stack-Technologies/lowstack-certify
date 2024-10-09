package main

import (
	"certify/internal/acme"
	"certify/internal/configuration"
	"fmt"
	"log"
)

func main() {
	log.Println("Low-Stack Certify is running!")

	// Initialize configuration and zones
	config := configuration.GetConfiguration()
	zones := acme.GetZones(config.ZonesPath)
	for _, zone := range zones {
		if err := acme.HandleZone(config, zone); err != nil {
			log.Fatal(fmt.Errorf("failed to handle zone: %w", err))
		}
	}
}
