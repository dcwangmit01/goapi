.PHONY: all certs hostdeps certs

NUM_CLIENTS := 2

all: certs

hostdeps:
	go get -u github.com/cloudflare/cfssl/cmd/...
	sudo pip install j2cli
	sudo apt-get -yq instal jq

certs:
	@# Initially create these files if they do not exist
	@#  cfssl print-defaults config > ca-cnf.json.j2
	@#  cfssl print-defaults csr > ca-csr.json.j2
	@# The above are hand edited and checked into git

	@# Apply the configuration to the templates
	j2 -f yaml ca-cnf.json.j2 config.yaml > ca-cnf.json
	j2 -f yaml ca-csr.json.j2 config.yaml > ca-csr.json

	@# Generate root CA certificate and private key
	mkdir -p ca
	if [ ! -f ./ca/ca-key.pem ]; then \
	  cfssl genkey -config=ca-cnf.json -profile=ca -initca ca-csr.json | cfssljson -bare ./ca/ca; \
	fi

	@# Generate client certificates signed by root CA
	mkdir -p clients
	for i in $$(seq -f "%02g" 01 $(NUM_CLIENTS)); do \
	  export CLIENT_NAME=client$$i; \
	  cat ca-csr.json | \
	    jq ".CN=\"$$CLIENT_NAME\"" | jq ".hosts=[\"$$CLIENT_NAME\"]" | \
	    cfssl gencert -config=ca-cnf.json -profile=client \
	      -ca=./ca/ca.pem -ca-key=./ca/ca-key.pem - \
	      | cfssljson -bare ./clients/$$CLIENT_NAME; \
	done

certcheck:
	for i in $$(seq -f "%02g" 01 $(NUM_CLIENTS)); do \
	  export CLIENT_NAME=client$$i; \
	  openssl x509 -text -in ./clients/$${CLIENT_NAME}.pem; \
	done
