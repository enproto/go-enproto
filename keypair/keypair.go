package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"os"
)

// KeyPair represents a pair of RSA keys (private and public).
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// GenerateRSAKeyPair generates a new RSA key pair.
func GenerateRSAKeyPair(size int) (*KeyPair, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	pubKey := &privKey.PublicKey

	return &KeyPair{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}, nil
}

// SaveKeyPair saves the RSA key pair to the specified files.
func SaveKeyPair(keyPair *KeyPair, privatePath, publicPath string) error {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair.PrivateKey)
	if err := os.WriteFile(privatePath, privKeyBytes, 0600); err != nil {
		return err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.PublicKey)
	if err != nil {
		return err
	}

	err = os.WriteFile(publicPath, pubKeyBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
