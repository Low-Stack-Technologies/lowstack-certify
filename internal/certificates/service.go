package certificates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	legoCertificate "github.com/go-acme/lego/v4/certificate"
	"os"
	"path"
	"time"
)

func GetExpirationDays(certificateDirectoryPath string) (int, error) {
	certificatePath := path.Join(certificateDirectoryPath, "cert1.pem")

	file, err := os.Open(certificatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to open certificate file: %w", err)
	}
	defer file.Close()

	certificateBytes, err := os.ReadFile(certificatePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read certificate file: %w", err)
	}

	block, _ := pem.Decode(certificateBytes)
	if block == nil {
		return 0, fmt.Errorf("failed to decode certificate file")
	}

	if block.Type != "CERTIFICATE" {
		return 0, fmt.Errorf("certificate file is not a certificate")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return 0, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return int(certificate.NotAfter.Sub(time.Now()).Hours() / 24), nil
}

func SaveCertificate(certificateDirectoryPath string, certificateResource *legoCertificate.Resource) error {
	fullChainPath := path.Join(certificateDirectoryPath, "fullchain.pem")
	if err := writeFileAndOverrideIfExists(fullChainPath, certificateResource.Certificate); err != nil {
		return fmt.Errorf("failed to write fullchain.pem: %w", err)
	}

	privateKeyPath := path.Join(certificateDirectoryPath, "privkey.pem")
	if err := writeFileAndOverrideIfExists(privateKeyPath, certificateResource.PrivateKey); err != nil {
		return fmt.Errorf("failed to write privkey.pem: %w", err)
	}

	return nil
}

func writeFileAndOverrideIfExists(path string, data []byte) error {
	if err := deleteFileIfExists(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func deleteFileIfExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}
