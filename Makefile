
.PHONY: vendor deps test clean

vendor:
	mkdir -p vendor/bin
	go build -o vendor/bin/gox vendor/github.com/mitchellh/gox/*.go
	go build -o vendor/bin/protoc-gen-go vendor/github.com/golang/protobuf/protoc-gen-go/*.go
	go build -o vendor/bin/protoc-gen-grpc-gateway vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/*.go
	go build -o vendor/bin/protoc-gen-swagger vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/*.go

deps:
	@# Link the parent of this current golang project directly into the GOPATH src
	@#   Golang needs to find the sources of this project in the GOPATH.
	@#   The parent directory is the org name 'dcwangmit01'
	if [ ! -L $$GOPATH/src/github.com/dcwangmit01 ]; then \
	  mkdir -p $$GOPATH/src/github.com; \
	  ln -s $$(readlink -f ../) $$GOPATH/src/github.com/dcwangmit01; \
	fi

	@# install the package dependencies in ./vendor
	glide install

gen:
	protoc \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	  entry/entry.proto
	protoc \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	  entry/entry.proto
	protoc -I/usr/local/include -I. \
	  -I . \
	  -I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	  entry/entry.proto

test:
	go test $(glide novendor)

clean:
	rm -rf vendor
