SHELL := /bin/bash

PATH := $(shell readlink -f ./bin/linux_amd64):$(shell readlink -f ./vendor/bin):$(PATH)
BIN_DIR := $(shell readlink -f ./bin)
PKG_DIR := $(shell readlink -f ./pkg)

# "ldflags" make go compile statically-linked binaries
GO_BUILD_FLAGS := -ldflags "-linkmode external -extldflags -static"


.PHONY: check deps vendor gen dist test clean mrclean

check:
	@# This is a check to make sure you run this makefile from within the GOPATH
	@#  Go requires building to be be run within GOPATH
	@if ! pwd | grep "$$GOPATH" > /dev/null; then \
	  echo "Cannot build unless within GOPATH $GOPATH"; \
	  echo "Please change to this current directory within GOPATH"; \
	fi

deps: check
	@# Link the parent of this current golang project directly into the GOPATH src
	@#   Golang needs to find the sources of this project in the GOPATH.
	@#   The parent directory is the org name 'dcwangmit01'
	if [ ! -L $$GOPATH/src/github.com/dcwangmit01 ]; then \
	  mkdir -p $$GOPATH/src/github.com; \
	  ln -s $$(readlink -f ../) $$GOPATH/src/github.com/dcwangmit01; \
	fi

	@# Install the package dependencies in ./vendor
	glide install

	@# install the arm cross compiler
	if ! which arm-linux-gnueabihf-gcc-5; then \
	  sudo apt-get -yq install gcc-5-arm-linux-gnueabihf; \
	fi

vendor: check deps
	@# Build tools on which this make system depends
	mkdir -p vendor/bin
	go build -o vendor/bin/protoc-gen-go vendor/github.com/golang/protobuf/protoc-gen-go/*.go
	go build -o vendor/bin/protoc-gen-grpc-gateway vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/*.go
	go build -o vendor/bin/protoc-gen-swagger vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/*.go

gen:
	@# Generate from the .proto file the GRPC definitons
	protoc \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	  entry-lib/entry.proto
	@# Generate from the .proto file the GRPC Gateway which proxies to JSON
	protoc \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	  entry-lib/entry.proto
	@# Generate from the .proto file the swagger definition
	protoc -I/usr/local/include -I. \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	  entry-lib/entry.proto

dist: check
	for GOOS in "linux"; do \
	  for GOARCH in "amd64" "arm"; do \
	    echo "Building $$GOOS-$$GOARCH"; \
	    export GOOS=$$GOOS; \
	    export GOARCH=$$GOARCH; \
	    if [ $$GOARCH = "arm" ]; then \
	      export GOARM=7; \
	      export CC=arm-linux-gnueabihf-gcc-5; \
	    fi; \
	    pushd ./entry-lib > /dev/null && \
	      go install -pkgdir="$(PKG_DIR)/$${GOOS}_$${GOARCH}" && \
	    popd > /dev/null ; \
	    pushd ./entry-server > /dev/null && \
	      go build $(GO_BUILD_FLAGS) \
	        -o "$(BIN_DIR)/$${GOOS}_$${GOARCH}/entry-server" \
	        -pkgdir="$(PKG_DIR)/$${GOOS}_$${GOARCH}" && \
	    popd > /dev/null ; \
	  done; \
	done

test:
	go test $(glide novendor)

clean:
	rm -rf bin pkg

mrclean: clean
	rm -rf vendor
