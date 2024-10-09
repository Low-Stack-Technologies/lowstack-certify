package provider_utils

import (
	"certify/internal/configuration"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

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

func GetACMEUser(configuration *configuration.Configuration, email string) (user *User, exists bool, err error) {
	// TODO: Check if user already exists

	// Generate private key for user
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		exists = false
		err = fmt.Errorf("failed to generate private key: %w", err)
		return
	}

	exists = false
	user = &User{
		Email: email,
		key:   privateKey,
	}

	return
}

func RegisterACMEUser(acmeClient *lego.Client, user *User) error {
	registrationResponse, err := acmeClient.Registration.Register(registration.RegisterOptions{
		TermsOfServiceAgreed: true,
	})
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	user.Registration = registrationResponse

	// TODO: Implement saving users to a file

	return nil
}
