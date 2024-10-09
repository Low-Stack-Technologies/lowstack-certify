package cloudflare

import (
	"certify/internal/acme/zone_configuration"
	"certify/internal/configuration"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	legoCertificate "github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	cloudflareChallenge "github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"os"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (p Provider) ObtainCertificate(configuration *configuration.Configuration, zoneConfiguration *zone_configuration.ZoneConfiguration) (*legoCertificate.Resource, error) {
	user, exists, err := getUser(configuration, zoneConfiguration.IdentityEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Create a new ACME config
	acmeConfig := lego.NewConfig(user)
	acmeConfig.CADirURL = configuration.CAURL
	acmeConfig.Certificate.KeyType = zoneConfiguration.KeyType

	// Create ACME client
	acmeClient, err := lego.NewClient(acmeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACME client: %w", err)
	}

	// Set the Cloudflare environment variables temporarily
	apiKey, ok := zoneConfiguration.ProviderOptions["api_token"]
	if !ok {
		return nil, fmt.Errorf("no cloudflare api token provided in configuration")
	}

	if err := os.Setenv("CLOUDFLARE_DNS_API_TOKEN", apiKey); err != nil {
		return nil, fmt.Errorf("failed to set Cloudflare API token: %w", err)
	}
	defer os.Setenv("CLOUDFLARE_DNS_API_TOKEN", "")

	// Cloudflare DNS Challenge
	challengeProvider, err := cloudflareChallenge.NewDNSProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare DNS provider: %w", err)
	}

	if err := acmeClient.Challenge.SetDNS01Provider(challengeProvider); err != nil {
		return nil, fmt.Errorf("failed to set Cloudflare DNS provider: %w", err)
	}

	// If the user does not exist, register them
	if !exists {
		reg, err := acmeClient.Registration.Register(registration.RegisterOptions{
			TermsOfServiceAgreed: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to register user: %w", err)
		}

		user.Registration = reg

		// TODO: Implement saving users to a file
	}

	// Obtain certificate
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

type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func getUser(configuration *configuration.Configuration, email string) (*User, bool, error) {
	// TODO: Implement saving and loading users from a file
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, false, fmt.Errorf("failed to generate private key: %w", err)
	}

	return &User{
		Email: email,
		key:   privateKey,
	}, false, nil
}
