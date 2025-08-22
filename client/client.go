package client

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/enproto/go-enproto/keypair"
)

// Client represents a client in the Enproto system.
type Client struct {
	publicKey *rsa.PublicKey
	client    *net.Conn
	config    *Config
}

// Config holds the configuration for the Enproto client.
type Config struct {
	ServerIP   string
	ServerPort int
}

// DefaultConfig returns the default configuration for the Enproto client.
func DefaultConfig() *Config {
	return &Config{
		ServerIP:   "127.0.0.1",
		ServerPort: 8080,
	}
}

// LoadClientFromFile loads the client's public key from the specified file.
func LoadClientFromFile(publicKeyPath string, config *Config) (*Client, error) {
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

	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		publicKey: pubKey,
		config:    config,
	}, nil
}

// LoadClient loads the client's public key from the specified files.
func LoadClient(pubKey keypair.PublicKey, config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		publicKey: pubKey,
		config:    config,
	}, nil
}

// Connect establishes a connection to the Enproto server.
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", c.config.ServerIP, c.config.ServerPort))
	if err != nil {
		return err
	}
	c.client = &conn
	return nil
}
