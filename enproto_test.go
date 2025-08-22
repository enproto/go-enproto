package enproto_test

import (
	"testing"

	"github.com/enproto/go-enproto/client"
	"github.com/enproto/go-enproto/keypair"
	"github.com/enproto/go-enproto/server"
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
		_, err := server.LoadServerFromFile("private.pem", "public.pem", nil)
		if err != nil {
			t.Fatalf("Failed to load server from file: %v", err)
		}

		_, err = client.LoadClientFromFile("public.pem", nil)
		if err != nil {
			t.Fatalf("Failed to load client from file: %v", err)
		}
	})

	t.Run("LoadFromKeyPair", func(t *testing.T) {
		kp, err := keypair.GenerateRSAKeyPair(2048)
		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}

		_, err = server.LoadServer(kp, &server.Config{TCPAddress: "127.0.0.11", TCPPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load server from key pair: %v", err)
		}

		_, err = client.LoadClient(kp.PublicKey, &client.Config{ServerIP: "127.0.0.11", ServerPort: 8080})
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

	var serv *server.Server
	var cli1 *client.Client
	var cli2 *client.Client
	var err error

	t.Run("Load", func(t *testing.T) {
		serv, err = server.LoadServerFromFile("private.pem", "public.pem", &server.Config{TCPAddress: "127.0.0.1", TCPPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load server from key pair: %v", err)
		}

		cli1, err = client.LoadClientFromFile("public.pem", &client.Config{ServerIP: "127.0.0.1", ServerPort: 8080})
		if err != nil {
			t.Fatalf("Failed to load client from key pair: %v", err)
		}

		cli2, err = client.LoadClientFromFile("public.pem", &client.Config{ServerIP: "127.0.0.1", ServerPort: 8080})
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
		recvs := serv.ReceiveClients()
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
}
