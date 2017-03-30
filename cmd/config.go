package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	certs "github.com/dcwangmit01/grpc-gw-poc/resources/certs"
)

var (
	keyPair       *tls.Certificate
	certPool      *x509.CertPool
	serverAddress string
	host          = "localhost"
	port          = 10080
)

func init() {

	var err error

	key, err := certs.Asset("insecure-key.pem")
	if err != nil {
		panic(err)
	}

	pem, err := certs.Asset("insecure.pem")
	if err != nil {
		panic(err)
	}

	pair, err := tls.X509KeyPair(pem, key)
	if err != nil {
		panic(err)
	}
	keyPair = &pair

	certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(pem)
	if !ok {
		panic("bad certs")
	}

	serverAddress = fmt.Sprintf("%s:%d", host, port)
}
