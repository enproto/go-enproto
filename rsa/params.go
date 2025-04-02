package rsa

type RSAParams struct {
	PrivateKeyPath string
	PublicKeyPath  string
	KeySize        int
}

func DefaultRSA() *RSAParams {
	return &RSAParams{
		PrivateKeyPath: "private_key.pem",
		PublicKeyPath:  "public_key.pem",
		KeySize:        2048,
	}
}
