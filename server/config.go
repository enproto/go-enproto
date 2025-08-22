package server

// Config holds the configuration for the Enproto server.
type Config struct {
	TCPAddress string
	TCPPort    int
}

// DefaultConfig returns the default configuration for the Enproto server.
func DefaultConfig() *Config {
	return &Config{
		TCPAddress: "0.0.0.0",
		TCPPort:    8080,
	}
}
