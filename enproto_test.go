package enproto_test

import (
	"testing"
	"time"

	"github.com/enproto/go-enproto/keypair"
	"github.com/enproto/go-enproto/service"
)

func TestEnprotoKeyGen(t *testing.T) {
	t.Run("KeyGenerate", func(t *testing.T) {
		kp, err := keypair.GenerateRSAKeyPair(2048)

		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}
		if err := keypair.SaveKeyPair(kp, "private.pem", "public.pem"); err != nil {
			t.Fatalf("Failed to save key pair: %v", err)
		}
	})

	t.Run("LoadFromFile", func(t *testing.T) {
		_, err := service.LoadServerFromFile("private.pem", "public.pem", nil)
		if err != nil {
			t.Fatalf("Failed to load server from file: %v", err)
		}

		_, err = service.LoadClientFromFile("public.pem", nil)
		if err != nil {
			t.Fatalf("Failed to load client from file: %v", err)
		}
	})

	t.Run("LoadFromKeyPair", func(t *testing.T) {
		kp, err := keypair.GenerateRSAKeyPair(2048)
		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}

		_, err = service.LoadServer(kp, &service.ServerConfig{TCPAddress: "127.0.0.11", TCPPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load server from key pair: %v", err)
		}

		_, err = service.LoadClient(kp, &service.ClientConfig{ServerIP: "127.0.0.11", ServerPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load client from key pair: %v", err)
		}
	})
}

func TestEnprotoService(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		kp, err := keypair.GenerateRSAKeyPair(2048)

		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}
		if err := keypair.SaveKeyPair(kp, "private.pem", "public.pem"); err != nil {
			t.Fatalf("Failed to save key pair: %v", err)
		}
	})

	var serv *service.Server
	var cli1 *service.Client
	var cli2 *service.Client
	var err error

	var recvs []*service.Client

	t.Run("Load", func(t *testing.T) {
		serv, err = service.LoadServerFromFile("private.pem", "public.pem", &service.ServerConfig{TCPAddress: "127.0.0.1", TCPPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load server from key pair: %v", err)
		}

		cli1, err = service.LoadClientFromFile("public.pem", &service.ClientConfig{ServerIP: "127.0.0.1", ServerPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load client from key pair: %v", err)
		}

		cli2, err = service.LoadClientFromFile("public.pem", &service.ClientConfig{ServerIP: "127.0.0.1", ServerPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load client from key pair: %v", err)
		}
	})

	t.Run("StartAndConnect", func(t *testing.T) {
		err = serv.Start()
		if err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}

		err = cli1.Connect()
		if err != nil {
			t.Fatalf("Failed to connect client: %v", err)
		}

		err = cli2.Connect()
		if err != nil {
			t.Fatalf("Failed to connect client: %v", err)
		}
	})

	t.Run("ReceiveClients", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)

		recvs = serv.ReceiveClients()
		if recvs == nil {
			t.Fatal("Failed to receive client")
		}

		if len(recvs) != 2 {
			t.Fatalf("Failed to receive all clients: %d", len(recvs))
		}
	})

	t.Run("Logging", func(t *testing.T) {
		for _, l := range serv.Logs() {
			t.Logf("Server log: %s", l)
		}
	})

	t.Run("ReadAndWrite", func(t *testing.T) {
		msg1 := []byte("Hello, Enproto!")
		msg2 := []byte("Goodbye, Enproto!")

		err = cli1.WriteSecure(msg1)
		if err != nil {
			t.Fatalf("Failed to write secure message: %v", err)
		}

		err = cli2.WriteSecure(msg1)
		if err != nil {
			t.Fatalf("Failed to write secure message: %v", err)
		}

		for _, c := range recvs {
			_, plaintext, err := c.ReadSecure()
			if err != nil {
				t.Fatalf("Failed to read secure message: %v", err)
			}
			if string(plaintext) != string(msg1) {
				t.Fatalf("Unexpected plaintext: %q", plaintext)
			}

			c.WriteSecure(msg2)
		}

		_, plain, err := cli1.ReadSecure()
		if err != nil {
			t.Fatalf("Failed to read secure message: %v", err)
		}
		if string(plain) != string(msg2) {
			t.Fatalf("Unexpected plaintext: %q", plain)
		}

		_, plain, err = cli2.ReadSecure()
		if err != nil {
			t.Fatalf("Failed to read secure message: %v", err)
		}
		if string(plain) != string(msg2) {
			t.Fatalf("Unexpected plaintext: %q", plain)
		}
	})
}
