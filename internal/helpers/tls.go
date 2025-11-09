// Package helpers implements utility functions for TLS certificate generation.
package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func GenerateSelfSignedTLSCert(host, certDir string) (certPEM, keyPEM string, err error) {
	// 0755 permissions to allow read and execute for everyone
	if err := os.MkdirAll(certDir, 0o755); err != nil {
		return "", "", fmt.Errorf("failed to create cert directory: %w", err)
	}

	// set file paths
	certPEM = filepath.Join(certDir, "cert.pem")
	keyPEM = filepath.Join(certDir, "key.pem")

	// check if already exist
	if _, err := os.Stat(certPEM); err == nil {
		if _, err := os.Stat(keyPEM); err == nil {
			log.Println("TLS certificate already exist")
			return certPEM, keyPEM, nil
		}
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"University of Technology"},
			Country:      []string{"IQ"},
			Locality:     []string{"baghdad"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Add host to certificate
	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = []net.IP{ip}
	} else {
		template.DNSNames = []string{host, "localhost"}
	}

	// DER : distinguished encoding rules
	// Create the binary certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %w", err)
	}

	// Write certificate file
	certOut, err := os.Create(certPEM)
	if err != nil {
		return "", "", fmt.Errorf("failed to open cert.pem for writing: %w", err)
	}
	defer func() {
		if err := certOut.Close(); err != nil {
			log.Printf("failed to close cert.pem: %v", err)
		}
	}()

	// encode cert to PEM format
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return "", "", fmt.Errorf("failed to write cert.pem: %w", err)
	}

	// Write private key file
	keyOut, err := os.OpenFile(keyPEM, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return "", "", fmt.Errorf("failed to open key.pem for writing: %w", err)
	}
	defer func() {
		if err := keyOut.Close(); err != nil {
			log.Printf("failed to close key.pem: %v", err)
		}
	}()

	// encode key to PEM format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER}); err != nil {
		return "", "", fmt.Errorf("failed to write key.pem: %w", err)
	}

	log.Printf("Generated self-signed certificate: %s, %s", certPEM, keyPEM)
	return certPEM, keyPEM, nil
}
