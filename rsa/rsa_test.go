package rsa_test

import (
	"testing"

	"github.com/enproto/go-enproto/rsa"
)

func TestRSA(t *testing.T) {
	rsa.GenerateKeypair(rsa.DefaultRSA())
}
