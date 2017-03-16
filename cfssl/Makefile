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
	j2 -f yaml client.ovpn.j2 config.yaml > client.ovpn

	@# Generate root CA certificate and private key
	mkdir -p ca
	if [ ! -f ./ca/ca-key.pem ]; then \
	  cfssl genkey -config=ca-cnf.json -profile=ca -initca ca-csr.json | cfssljson -bare ./ca/ca; \
	fi

	@# Generate single server and many client certificates signed by root CA
	mkdir -p clients
	for i in server $$(seq -f "%02g" 01 $(NUM_CLIENTS) | awk '{print "client" $$0}'); do \
	  if [ ! -f ./clients/$$i.pem ]; then \
	    cat ca-csr.json | \
	      jq ".CN=\"$$i\"" | jq ".hosts=[\"$$i\"]" | \
	      cfssl gencert -config=ca-cnf.json -profile=client \
	        -ca=./ca/ca.pem -ca-key=./ca/ca-key.pem - \
	        | cfssljson -bare ./clients/$$i; \
	  fi; \
	done

	@# Generate OpenVPN configs per client
	for i in $$(seq -f "%02g" 01 $(NUM_CLIENTS) | awk '{print "client" $$0}'); do \
	  export C=clients/$$i; \
	  cp client.ovpn $$C.ovpn; \
	  echo "" >> $$C.ovpn; \
	  echo "<ca>" >> $$C.ovpn; \
	  cat ./ca/ca.pem >> $$C.ovpn; \
	  echo "</ca>" >> $$C.ovpn; \
	  echo "" >> $$C.ovpn; \
	  echo "<cert>" >> $$C.ovpn; \
	  cat $$C.pem >> $$C.ovpn; \
	  echo "</cert>" >> $$C.ovpn; \
	  echo "" >> $$C.ovpn; \
	  echo "<key>" >> $$C.ovpn; \
	  cat $$C-key.pem >> $$C.ovpn; \
	  echo "</key>" >> $$C.ovpn; \
	  echo "" >> $$C.ovpn; \
	  echo "<dh>" >> $$C.ovpn; \
	  cat ./ca/dh4096.pem >> $$C.ovpn; \
	  echo "</dh>" >> $$C.ovpn; \
	done

	@# Generate Diffie-Hellman Params
	if [ ! -f ./ca/dh4096.pem ]; then \
	  openssl dhparam -out ./ca/dh4096.pem 4096; \
	fi

mrclean:
	rm -rf ./ca ./clients *.json *.ovpn


certcheck:
	for i in $$(seq -f "%02g" 01 $(NUM_CLIENTS)); do \
	  export CLIENT_NAME=client$$i; \
	  openssl x509 -text -in ./clients/$$CLIENT_NAME.pem; \
	done
