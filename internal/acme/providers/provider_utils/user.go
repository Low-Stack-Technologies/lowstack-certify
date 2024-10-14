package provider_utils

import (
	"certify/internal/configuration"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"log"
	"os"
	"path"
	"strings"
)

type AcmeUserFile struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	KeyPem       string                 `json:"key_pem"`
}

type User struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
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

func (u *User) SaveToFile() error {
	usersPath := getACMEUsersPath()
	if err := makeDirectoryIfNotExists(usersPath); err != nil {
		return fmt.Errorf("failed to create users directory: %w", err)
	}

	filePath := path.Join(usersPath, emailToFileSafeString(u.Email))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create user file: %w", err)
	}
	defer file.Close()

	derKey, err := x509.MarshalECPrivateKey(u.key.(*ecdsa.PrivateKey))
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: derKey})

	acmeUserFileContent := AcmeUserFile{
		Email:        u.Email,
		Registration: u.Registration,
		KeyPem:       string(pemKey),
	}

	if err := json.NewEncoder(file).Encode(acmeUserFileContent); err != nil {
		return fmt.Errorf("failed to encode user to JSON: %w", err)
	}

	if err := os.Chmod(filePath, 0600); err != nil {
		return fmt.Errorf("failed to set user file permissions: %w", err)
	}

	return nil
}

func LoadACMEUser(configuration *configuration.Configuration, email string) (user *User, err error) {
	usersPath := getACMEUsersPath()
	filePath := path.Join(usersPath, emailToFileSafeString(email))

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open user file: %w", err)
	}
	defer file.Close()

	var acmeUserFile AcmeUserFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&acmeUserFile); err != nil {
		return nil, fmt.Errorf("failed to decode user file: %w", err)
	}

	block, _ := pem.Decode([]byte(acmeUserFile.KeyPem))
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM block: %w", err)
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	user = &User{
		Email:        acmeUserFile.Email,
		Registration: acmeUserFile.Registration,
		key:          privateKey,
	}

	return user, nil
}

func GetACMEUser(configuration *configuration.Configuration, email string) (user *User, exists bool, err error) {
	// Try to load user from file
	user, err = LoadACMEUser(configuration, email)
	if err != nil {
		log.Printf("Failed to load user from file: %s", err)
	}

	// Return if user was loaded
	if user != nil {
		log.Printf("User %s was loaded from file", email)
		return user, true, nil
	}

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

	if err := user.SaveToFile(); err != nil {
		log.Printf("Failed to save user to file: %s", err)
	}

	return nil
}

func getACMEUsersPath() string {
	configurationDirectory := path.Dir(configuration.GetConfigurationPath())
	return path.Join(configurationDirectory, "acme_users")
}

func emailToFileSafeString(email string) string {
	return strings.ReplaceAll(email, "@", "_")
}
