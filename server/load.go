package server

import (
	"crypto/rsa"
	"crypto/x509"
	"os"

	"github.com/enproto/go-enproto/internal/util"
	"github.com/enproto/go-enproto/keypair"
)

// LoadServerFromFile loads the server's key pair from the specified files.
func LoadServerFromFile(privateKeyPath string, publicKeyPath string, config *Config) (*Server, error) {
	privKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	keypair := &keypair.KeyPair{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}

	if config == nil {
		config = DefaultConfig()
	}

	logChan := make(chan string, 100)
	logChan <- util.MakeLog("Server loaded from files")

	return &Server{
		keyPair:    keypair,
		config:     config,
		clientChan: make(chan *Client),
		errorChan:  make(chan error),
		logChan:    logChan,
	}, nil
}

// LoadServer loads the server's key pair from the specified files.
func LoadServer(keyPair *keypair.KeyPair, config *Config) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}

	logChan := make(chan string, 100)
	logChan <- util.MakeLog("Server loaded from key pair")

	return &Server{
		keyPair:    keyPair,
		config:     config,
		clientChan: make(chan *Client),
		errorChan:  make(chan error),
		logChan:    logChan,
	}, nil
}
