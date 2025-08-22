package service

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

// Client represents a client in the Enproto system.
type Client struct {
	publicKey *rsa.PublicKey
	conn      *net.TCPConn
	config    *ClientConfig
	isLoaded  bool
	aesKey    []byte
	txSeq     uint64 // transmit sequence (my outgoing)
	rxSeq     uint64 // last accepted receive sequence (peer's incoming)
	aead      cipher.AEAD
}

// Connect establishes a connection to the Enproto server.
func (c *Client) Connect() error {
	if !c.isLoaded {
		return fmt.Errorf("client is not loaded")
	}

	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: net.ParseIP(c.config.ServerIP), Port: c.config.ServerPort})
	if err != nil {
		return err
	}

	c.conn = conn
	c.aesKey = make([]byte, 32)

	_, err = rand.Read(c.aesKey)
	if err != nil {
		return err
	}

	enc, err := rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, c.aesKey)
	if err != nil {
		return err
	}

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(enc)))

	bw := bufio.NewWriter(c.conn)

	_, err = bw.Write(lenBuf[:])
	if err != nil {
		return err
	}

	_, err = bw.Write(enc)
	if err != nil {
		return err
	}

	err = bw.Flush()
	if err != nil {
		return err
	}

	return nil
}

const (
	wireVersion     = byte(1)    // protocol version (bump if wire layout changes)
	headerLen       = 1 + 8 + 12 // version(1) + seq(8) + nonce(12)
	maxCipherPacket = 8 << 20    // 8 MiB sane upper bound for ciphertext (tune as needed)
)

// buildAEAD constructs an AES-GCM AEAD from c.aesKey.
func (c *Client) buildAEAD() (cipher.AEAD, error) {
	if c.aead != nil {
		return c.aead, nil
	}

	if len(c.aesKey) == 0 {
		return nil, errors.New("aes key not set")
	}
	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}
	// gcm.NonceSize() is expected to be 12
	c.aead = gcm

	return gcm, nil
}

// ---- Secure Write ----
// Frame: [4B length] [1B version] [8B seq] [12B nonce] [ciphertext||tag]
// AAD: version||seq
func (c *Client) WriteSecure(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	if len(p) > maxCipherPacket {
		return fmt.Errorf("plaintext too large: %d", len(p))
	}
	gcm, err := c.buildAEAD()
	if err != nil {
		return err
	}

	// Sequence strictly increases per outgoing message.
	seq := atomic.AddUint64(&c.txSeq, 1) - 1

	// Build header
	hdr := make([]byte, headerLen)
	hdr[0] = wireVersion
	binary.BigEndian.PutUint64(hdr[1:9], seq)

	// Nonce: 12 bytes; must be unique per (key, message).
	// In production, prepend a session/epoch + direction prefix instead of pure random.
	nonce := hdr[9:21]
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("nonce: %w", err)
	}

	aad := hdr[:1+8] // version||seq
	ct := gcm.Seal(nil, nonce, p, aad)

	total := len(hdr) + len(ct)
	if total > headerLen+maxCipherPacket {
		return fmt.Errorf("ciphertext too large: %d", total)
	}

	// Length prefix (big-endian)
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(total))

	if err := writeFull(c.conn, lenBuf[:]); err != nil {
		return fmt.Errorf("write len: %w", err)
	}
	if err := writeFull(c.conn, hdr); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	if err := writeFull(c.conn, ct); err != nil {
		return fmt.Errorf("write body: %w", err)
	}
	return nil
}

func writeFull(w io.Writer, b []byte) error {
	n, err := w.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return io.ErrShortWrite
	}
	return nil
}

// ---- Secure Read ----
// Returns (seq, plaintext).
func (c *Client) ReadSecure() (uint64, []byte, error) {
	gcm, err := c.buildAEAD()
	if err != nil {
		return 0, nil, err
	}

	// 1) Read length
	var lenBuf [4]byte
	if _, err := io.ReadFull(c.conn, lenBuf[:]); err != nil {
		return 0, nil, fmt.Errorf("read len: %w", err)
	}
	total := int(binary.BigEndian.Uint32(lenBuf[:]))
	if total < headerLen {
		return 0, nil, fmt.Errorf("invalid frame: short %d", total)
	}
	if total > headerLen+maxCipherPacket {
		return 0, nil, fmt.Errorf("invalid frame: large %d", total)
	}

	// 2) Read frame
	frame := make([]byte, total)
	if _, err := io.ReadFull(c.conn, frame); err != nil {
		return 0, nil, fmt.Errorf("read frame: %w", err)
	}

	// 3) Parse & decrypt
	hdr := frame[:headerLen]
	if hdr[0] != wireVersion {
		return 0, nil, fmt.Errorf("unsupported version: %d", hdr[0])
	}
	seq := binary.BigEndian.Uint64(hdr[1:9])
	nonce := hdr[9:21]
	aad := hdr[:1+8]
	ct := frame[headerLen:]

	pt, err := gcm.Open(nil, nonce, ct, aad)
	if err != nil {
		return 0, nil, fmt.Errorf("decrypt: %w", err) // integrity/auth failure
	}

	// 4) Anti-replay / ordering policy
	// Strict monotonic: accept only seq == c.rxSeq+1
	last := atomic.LoadUint64(&c.rxSeq)
	if seq != last+1 {
		// If you prefer to drop silently:
		// return 0, nil, fmt.Errorf("out-of-order or replay: got %d want %d", seq, last+1)

		// Alternatively, implement sliding-window:
		// if !c.acceptInWindow(seq) { return 0, nil, fmt.Errorf("replay old seq %d", seq) }
		// else continue
		return 0, nil, fmt.Errorf("out-of-order or replay: got %d want %d", seq, last+1)
	}

	atomic.StoreUint64(&c.rxSeq, seq)
	return seq, pt, nil
}
