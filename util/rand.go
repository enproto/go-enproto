package util

import "crypto/rand"

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
