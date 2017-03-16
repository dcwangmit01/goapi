.PHONY: all certs hostdeps certs


# There isn't much different between the servers, peers, and client
#   certs below.  They are just named differently to conveniently
#   identify different uses.
# Network servers
NUM_SERVERS := 5
# Peers are site to site
NUM_PEERS := 5
# OVPN Client certs
NUM_CLIENTS := 20

# Keep CA files separate than certs for safety
CA_DIR := $(shell readlink -f ./ca)
CERTS_DIR := $(shell readlink -f ./certs)

# Ensure the dirs above exist on a clean checkout
$(shell mkdir -p $(CA_DIR) $(CERTS_DIR))


.PHONY: hostdeps templates certs

all: templates ovpn

hostdeps:
	@# On a linux host, install the following
	go get -u github.com/cloudflare/cfssl/cmd/...
	sudo pip install j2cli
	sudo apt-get -yq instal jq


templates:
	@# Initially create these files if they do not exist
	@#  cfssl print-defaults config > ca-cnf.json.j2
	@#  cfssl print-defaults csr > ca-csr.json.j2
	@# The above are hand edited and checked into git

	@# Apply the configuration to the templates
	j2 -f yaml ca-cnf.json.j2 config.yaml > ca-cnf.json
	j2 -f yaml ca-csr.json.j2 config.yaml > ca-csr.json
	j2 -f yaml client.ovpn.j2 config.yaml > client.ovpn


$(CA_DIR)/dh4096.pem:
	@# Generate Diffie-Hellman Params
	openssl dhparam -out $(CA_DIR)/dh4096.pem 4096


$(CA_DIR)/ca-key.pem:
	@# Generate root CA certificate and private key
	cfssl genkey -config=ca-cnf.json -profile=ca -initca ca-csr.json | cfssljson -bare $(CA_DIR)/ca


certs: $(CA_DIR)/ca-key.pem
	@# Generate all certs signed by root CA
	for i in $$(seq -f "%02g" 01 $(NUM_SERVERS) | awk '{print "server" $$0}') \
	         $$(seq -f "%02g" 01 $(NUM_CLIENTS) | awk '{print "client" $$0}') \
	         $$(seq -f "%02g" 01 $(NUM_PEERS) | awk '{print "peer" $$0}') ; do \
	  if [ ! -f $(CERTS_DIR)/$$i.pem ]; then \
	    cat ca-csr.json | \
	      jq ".CN=\"$$i\"" | jq ".hosts=[\"$$i\"]" | \
	      cfssl gencert -config=ca-cnf.json -profile=client \
	        -ca=$(CA_DIR)/ca.pem -ca-key=$(CA_DIR)/ca-key.pem - | \
	      cfssljson -bare $(CERTS_DIR)/$$i; \
	  fi; \
	done

ovpn: certs $(CA_DIR)/dh4096.pem
	@# Generate OpenVPN configs per client
	for i in $$(seq -f "%02g" 01 $(NUM_CLIENTS) | awk '{print "client" $$0}'); do \
	  export C=$(CERTS_DIR)/$$i; \
	  cp client.ovpn $$C.ovpn; \
	  echo "" >> $$C.ovpn; \
	  echo "<ca>" >> $$C.ovpn; \
	  cat $(CA_DIR)/ca.pem >> $$C.ovpn; \
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
	  cat $(CA_DIR)/dh4096.pem >> $$C.ovpn; \
	  echo "</dh>" >> $$C.ovpn; \
	done


mrclean:
	rm -rf $(CA_DIR) $(CERTS_DIR) *.json *.ovpn


certcheck:
	for i in $$(seq -f "%02g" 01 $(NUM_CLIENTS) | awk '{print "client" $$0}'); do \
	  export C=clients/$$i; \
	  openssl x509 -text -in $$C.pem; \
	done
