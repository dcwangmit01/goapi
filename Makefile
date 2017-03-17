SHELL := /bin/bash
PACKAGE := github.com/dcwangmit01/grpc-gw-poc

BIN_NAME := $(basename $(PACKAGE))

# Modify the current path to use locally built tools
PATH := $(shell readlink -f ./bin/linux_amd64):$(shell readlink -f ./vendor/bin):$(PATH)
CWD := $(shell readlink -f ./)

GO_DIR := $(GOPATH)/src/$(PACKAGE)

BIN_DIR := $(shell readlink -f ./bin)
PKG_DIR := $(shell readlink -f ./pkg)
BUILD_DIR := $(shell readlink -f ./.build)
CACHE_DIR := $(shell readlink -f ./.cache)
CERTS_DIR := $(shell readlink -f ./cfssl/certs)
RESOURCE_DIR := $(shell readlink -f ./resources)
# Ensure the dirs above exist on a clean checkout
$(shell mkdir -p $(BIN_DIR) $(PKG_DIR) $(BUILD_DIR) $(CACHE_DIR) $(CERTS_DIR) $(RESOURCE_DIR))

# Ensure go compiles statically-linked binaries with "ldflags"
GO_BUILD_FLAGS := -ldflags "-linkmode external -extldflags -static"

# Swagger version to package and deploy
SWAGGER_UI_VERSION := 2.2.8


.PHONY: hostdeps check vendor code_gen resource_gen cert_gen dist_all

all: dist

hostdeps:
	@# install the arm cross compiler
	@if ! which arm-linux-gnueabihf-gcc-5 > /dev/null; then \
	  sudo apt-get -yq install gcc-5-arm-linux-gnueabihf; \
	fi

	@# Link this current golang project directly into the GOPATH/src
	@#   Golang needs to find the sources of this project in the GOPATH.
	@#   The parent directory is the org name 'dcwangmit01'
	@if [ ! -L $(GO_DIR) ]; then \
	  mkdir -p $(shell dirname $(GO_DIR)); \
	  ln -s $(CWD) $(GO_DIR); \
	fi

check: hostdeps
	@# This is a check to make sure you run this makefile from within the GOPATH
	@#  Go requires building to be be run within GOPATH
	@if ! pwd | grep "$$GOPATH" > /dev/null; then \
	  echo "Cannot build unless within GOPATH workspace at $$GOPATH"; \
	  echo "Please change to this current directory within GOPATH"; \
	  echo ""; \
	  echo "  cd $(GO_DIR)"; \
	  echo ""; \
	  exit 1; \
	fi

vendor: check
	@# Work around "directory not empty" bug on second glide up/install
	rm -rf vendor

	@# Install the package dependencies in ./vendor
	glide install

	@# Build the tools on which this project build system depends
	mkdir -p vendor/bin
	go build -o vendor/bin/protoc-gen-go vendor/github.com/golang/protobuf/protoc-gen-go/*.go
	go build -o vendor/bin/protoc-gen-grpc-gateway vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/*.go
	go build -o vendor/bin/protoc-gen-swagger vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/*.go
	go build -o vendor/bin/go-bindata vendor/github.com/jteeuwen/go-bindata/go-bindata/*.go
	go build -o vendor/bin/go-bindata-assetfs vendor/github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs/*.go

code_gen: check
	@# Generate the GRPC definitons from the .proto file
	protoc \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	  entry/entry.proto

	@# Generate the GRPC Gateway which proxies to JSON from the .proto file
	protoc \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	  entry/entry.proto

	@# Generate the swagger definition from the .proto file
	protoc -I/usr/local/include -I. \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	  entry/entry.proto

resource_gen: check
	@# Download and extract the swagger-ui release
	export SWAGGER_UI_TGZ=$(CACHE_DIR)/swagger-ui-$(SWAGGER_UI_VERSION).tar.gz; \
	if [ ! -f $$SWAGGER_UI_TGZ ]; then \
	  curl -fsSL https://github.com/swagger-api/swagger-ui/archive/v$(SWAGGER_UI_VERSION).tar.gz > \
	    $$SWAGGER_UI_TGZ; \
	fi; \
	if [ ! -d $(CACHE_DIR)/swagger-ui-$(SWAGGER_UI_VERSION) ]; then \
	  tar xzf $$SWAGGER_UI_TGZ -C $(BUILD_DIR); \
	fi

	@# Generate the swagger-ui directory as a golang file
	@# Ignore the warning about "Cannot read bindata.go open bindata.go: no such file or directory"
	mkdir -p $(RESOURCE_DIR)/swagger/ui
	cd $(BUILD_DIR)/swagger-ui-$(SWAGGER_UI_VERSION)/dist && \
	  go-bindata-assetfs -o $(RESOURCE_DIR)/swagger/ui/ui.go -pkg ui 2>/dev/null ./... || true

	@# Generate the swagger.json file as a golang file
	mkdir -p $(RESOURCE_DIR)/swagger/files
	mkdir -p $(BUILD_DIR)/swagger/files
	cp -f $(CWD)/entry/entry.swagger.json $(BUILD_DIR)/swagger/files/swagger.json
	cd $(BUILD_DIR)/swagger/files && \
	  go-bindata-assetfs -o $(RESOURCE_DIR)/swagger/files/files.go -pkg files 2>/dev/null ./... || true

cert_gen:
	@# Call the Makefile in the subdir
	make -C cfssl

	@# Generate the certs/entry/ directory as a golang file
	@# Ignore the warning about "Cannot read bindata.go open bindata.go: no such file or directory"
	mkdir -p $(RESOURCE_DIR)/certs
	cd $(CERTS_DIR) && \
	  go-bindata-assetfs -o $(RESOURCE_DIR)/certs/certs.go -pkg certs ./... 2>/dev/null || true

dist: check
	export GOOS=linux; \
	export GOARCH=amd64; \
	go build $(GO_BUILD_FLAGS) \
	  -o "$(BIN_DIR)/$${GOOS}_$${GOARCH}/$(BIN_NAME)"

dist_all: hostdeps check vendor code_gen resource_gen cert_gen dist_all
	for GOOS in "linux"; do \
	  for GOARCH in "amd64" "arm"; do \
	    echo "Building $$GOOS-$$GOARCH"; \
	    export GOOS=$$GOOS; \
	    export GOARCH=$$GOARCH; \
	    if [ $$GOARCH = "arm" ]; then \
	      export GOARM=7; \
	      export CC=arm-linux-gnueabihf-gcc-5; \
	    fi; \
	    go build $(GO_BUILD_FLAGS) \
	      -o "$(BIN_DIR)/$${GOOS}_$${GOARCH}/$(BIN_NAME)"; \
	  done; \
	done

test:
	go test $(glide novendor)

clean:
	rm -rf bin pkg .build

mrclean: clean
	rm -rf vendor


notes:
	@### Notes. The following can be used to build a lib file
	@# cd ./entry && \
	@#   go install -pkgdir="$(PKG_DIR)/$${GOOS}_$${GOARCH}"
