package server

import (
	"fmt"
	"net"

	"github.com/enproto/go-enproto/internal/util"
	"github.com/enproto/go-enproto/keypair"
)

// Server represents a server in the Enproto system.
type Server struct {
	keyPair     *keypair.KeyPair
	tcpListener *net.Listener
	config      *Config
	clientChan  chan *Client
	errorChan   chan error
	logChan     chan string
}

// Start starts the Enproto server.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", s.config.TCPAddress, s.config.TCPPort))
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				s.errorChan <- err
			}

			go s.handleConnection(&conn)
		}
	}()

	return nil
}

// handleConnection handles a new client connection.
func (s *Server) handleConnection(conn *net.Conn) {
	client := &Client{
		conn: conn,
	}

	s.logChan <- util.MakeLog("New client connected: %s", (*conn).RemoteAddr().String())
	s.clientChan <- client
}

// ReceiveClients retrieves a client from the server's client channel.
func (s *Server) ReceiveClients() []*Client {
	var clients []*Client

	for {
		select {
		case client := <-s.clientChan:
			clients = append(clients, client)
		default:
			return clients
		}
	}
}

// Logs retrieves a log message from the server's log channel.
func (s *Server) Logs() []string {
	var logs []string

	for {
		select {
		case log := <-s.logChan:
			logs = append(logs, log)
		default:
			return logs
		}
	}
}
