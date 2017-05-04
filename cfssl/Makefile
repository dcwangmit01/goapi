.PHONY: all certs hostdeps certs

CERTS_DIR := $(shell readlink -f ./certs)

# Ensure the dirs above exist on a clean checkout
$(shell mkdir -p $(CERTS_DIR))

.PHONY: hostdeps certs

all: certs

hostdeps:
	@# On a linux host, install the following
	go get -u github.com/cloudflare/cfssl/cmd/...

certs: cfssl/certs/insecure-key.pem cfssl/certs/insecure.pem

cfssl/certs/insecure-key.pem cfssl/certs/insecure.pem:
	cfssl selfsign localhost cfssl-csr.json | cfssljson -bare $(CERTS_DIR)/insecure

mrclean:
	rm -rf $(CA_DIR) $(CERTS_DIR)

certcheck:
	openssl x509 -text -in $(CERTS_DIR)/insecure.pem
