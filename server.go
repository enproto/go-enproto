package enproto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Server struct {
	privKey *rsa.PrivateKey
}

func NewServer(privKeyPath string) (*Server, error) {
	data, err := os.ReadFile(privKeyPath)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid PEM block")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return &Server{
		privKey: key,
	}, nil
}
