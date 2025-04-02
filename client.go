package enproto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Client struct {
	pubKey *rsa.PublicKey
}

func NewClient(pubKeyPath string) (*Client, error) {
	data, err := os.ReadFile(pubKeyPath)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}

	return &Client{
		pubKey: rsaPub,
	}, nil
}
