package service

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/enproto/go-enproto/internal/util"
	"github.com/enproto/go-enproto/keypair"
)

// Server represents a server in the Enproto system.
type Server struct {
	keyPair    *keypair.KeyPair
	listener   *net.TCPListener
	config     *ServerConfig
	clientChan chan *Client
	logChan    chan string
}

// Start starts the Enproto server.
func (s *Server) Start() error {
	ln, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.ParseIP(s.config.TCPAddress), Port: s.config.TCPPort})
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				s.logChan <- fmt.Sprintf("Error accepting connection: %v", err)
				continue
			}

			go s.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection handles a new client connection.
func (s *Server) handleConnection(conn *net.TCPConn) {
	client := &Client{
		conn:     conn,
		isLoaded: false,
		txSeq:    1,
		rxSeq:    0,
	}

	s.logChan <- util.MakeLog("New client connected: %s", conn.RemoteAddr().String())

	br := bufio.NewReader(conn)

	var lenBuf [4]byte

	_, err := io.ReadFull(br, lenBuf[:])
	if err != nil {
		s.logChan <- util.MakeLog("Error reading length prefix: %v", err)
		return
	}

	size := binary.BigEndian.Uint32(lenBuf[:])
	keyBytes := (s.keyPair.PrivateKey.N.BitLen() + 7) / 8

	if int(size) != keyBytes {
		s.logChan <- util.MakeLog("Invalid key size: got %d, want %d", size, keyBytes)
		return
	}

	buf := make([]byte, size)

	_, err = io.ReadFull(br, buf)
	if err != nil {
		s.logChan <- util.MakeLog("Error reading AES key: %v", err)
		return
	}

	dec, err := rsa.DecryptPKCS1v15(rand.Reader, s.keyPair.PrivateKey, buf)
	if err != nil {
		s.logChan <- util.MakeLog("Error decrypting AES key: %v", err)
		return
	}

	client.aesKey = dec

	s.logChan <- util.MakeLog("Client %s connected with AES key: %x...", conn.RemoteAddr().String(), client.aesKey[:4])
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
