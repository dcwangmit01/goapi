package certs

import (
	"crypto/tls"
	"crypto/x509"
)

var (
	KeyPair  *tls.Certificate // contains private key and public cert
	CertPool *x509.CertPool   // a certpool is a set of certificates.  Only contains public cert
)

func init() {

	var err error

	key, err := Asset("insecure-key.pem")
	if err != nil {
		panic(err)
	}

	pem, err := Asset("insecure.pem")
	if err != nil {
		panic(err)
	}

	pair, err := tls.X509KeyPair(pem, key)
	if err != nil {
		panic(err)
	}

	KeyPair = &pair

	CertPool = x509.NewCertPool()
	ok := CertPool.AppendCertsFromPEM(pem)
	if !ok {
		panic("bad certs")
	}
}
