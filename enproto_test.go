package enproto_test

import (
	"log"
	"testing"

	"github.com/enproto/go-enproto"
	"github.com/enproto/go-enproto/rsa"
)

func TestRSA(t *testing.T) {
	rsa.GenerateKeypair(rsa.DefaultRSA())

	_, err := enproto.NewServer(rsa.DefaultRSA().PrivateKeyPath)

	if err != nil {
		log.Fatalf("Do not generate server")
	}

	_, err = enproto.NewClient(rsa.DefaultRSA().PublicKeyPath)

	if err != nil {
		log.Fatalf("Do not generate client")
	}
}
