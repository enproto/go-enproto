package service

// ServerConfig holds the configuration for the Enproto server.
type ServerConfig struct {
	TCPAddress string
	TCPPort    int
}

// DefaultServerConfig returns the default configuration for the Enproto server.
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		TCPAddress: "0.0.0.0",
		TCPPort:    8080,
	}
}

// ClientConfig holds the configuration for the Enproto client.
type ClientConfig struct {
	ServerIP   string
	ServerPort int
}

// DefaultClientConfig returns the default configuration for the Enproto client.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ServerIP:   "127.0.0.1",
		ServerPort: 8080,
	}
}
