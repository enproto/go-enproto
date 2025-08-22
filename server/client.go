package server

import "net"

// Client represents a client in the Enproto system.
type Client struct {
	conn *net.Conn
}
