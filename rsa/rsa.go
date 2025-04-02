package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type RSAParams struct {
	PrivateKeyPath string
	PublicKeyPath  string
	KeySize        int
}

func DefaultRSA() *RSAParams {
	return &RSAParams{
		PrivateKeyPath: "private_key.pem",
		PublicKeyPath:  "public_key.pem",
		KeySize:        2048,
	}
}

const ()

func GenerateKeypair(params *RSAParams) {
	if params.KeySize != 1024 && params.KeySize != 2048 && params.KeySize != 4096 {
		log.Fatalf("Invalid key size: The key size can not %d", params.KeySize)
	}

	// Generate RSA Key Pair
	privateKey, err := rsa.GenerateKey(rand.Reader, params.KeySize)

	if err != nil {
		log.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Encode and save private key in PEM format
	err = savePrivateKey(params.PrivateKeyPath, privateKey)

	if err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	log.Printf("Private key saved to \"%s\"\n", params.PrivateKeyPath)

	// Encode and save public key in PEM format
	err = savePublicKey(params.PublicKeyPath, &privateKey.PublicKey)

	if err != nil {
		log.Fatalf("Failed to save public key: %v", err)
	}

	log.Printf("Public key saved to \"%s\"\n", params.PublicKeyPath)
}

func savePrivateKey(filename string, key *rsa.PrivateKey) error {
	// Convert to PKCS#1 ASN.1 DER encoded form
	privDER := x509.MarshalPKCS1PrivateKey(key)

	// Encode to PEM
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})

	// Write to file
	return os.WriteFile(filename, privPEM, 0600)
}

func savePublicKey(filename string, pub *rsa.PublicKey) error {
	// Convert to PKIX ASN.1 DER encoded form
	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	// Encode to PEM
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	})

	// Write to file
	return os.WriteFile(filename, pubPEM, 0644)
}
