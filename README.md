# goapi

This repo is a forkable base project for implementing golang API servers.

It uses [GRPC](http://www.grpc.io/) and
[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) instead of
typical API toolkits such as [Gin](https://github.com/gin-gonic/gin) or
[Gorilla](http://www.gorillatoolkit.org/).  This has the benefits of
code-generation.

This project provides:

* An API server implementation which services requests over
  [GRPC](http://www.grpc.io/) [protocol
  buffers](https://developers.google.com/protocol-buffers/) as well as JSON
  REST.  All JSON REST requests are proxied over
  [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to the GRPC
  server.

* A [cobra](https://github.com/spf13/cobra) CLI tool with subcommands that can
  launch the API server as well as drive the API as a CLI client.  The CLI tool
  stores its settings in [viper](https://github.com/spf13/viper)
  configuration.

* A Makefile which generates code resources including:
    * GRPC code files including protobuf, grpc, grpc-gateway, swagger
    * SSL certs using [cfssl](https://github.com/cloudflare/cfssl)
    * Swagger asset files with embedded swagger-ui, to serve the swagger file
    * Cross-compiling statically linked Linux binaries for amd64 and armv7.

* A test framework using [ginkgo](https://github.com/onsi/ginkgo) and [gomega
  matchers](https://github.com/onsi/gomega)

* A semantic version-based vendoring system using
  [glide](https://github.com/Masterminds/glide)

* Using a vendored golang package management framework (glide)

* Provides an `/auth` endpoint which returns a JWT token when presented with
  valid credentials.  This authentication system is based on
  [JWT](github.com/dgrijalva/jwt-go)

* An [Role Based Access Control](https://github.com/mikespook/gorbac)
  authorization system based on

* A structured logging system using
  [logrus](https://github.com/sirupsen/logrus)

TODO

* The tools to code-generate the API framework from services and endpoints
  described in a `.proto` service description file.
* An example key/value api storing values to sqlite.
* A user database




## Building

```
make
```

## Running

```
# Show Usage
./bin/linux_amd64/goapi

# Start the GRPC+JSON server on port 10080
./bin/linux_amd64/goapi serve

### Use the GRPC CLI Client

# Obtain an Auth Token
./bin/linux_amd64/goapi auth admin password

# Create a Key
./bin/linux_amd64/goapi keyval create mykey myval

# Read a Key
./bin/linux_amd64/goapi keyval read mykey

# Update a Key
./bin/linux_amd64/goapi keyval create mykey myval2

# Delete a Key
./bin/linux_amd64/goapi keyval delete mykey

### Use Curl

# Obtain an Auth Token
curl -vvv -X POST -k https://localhost:10080/v1/auth -H "Content-Type: text/plain" -d '{"email": "admin", "password": "password"}'

# Create a Key
curl -vvv -X PUT -k https://localhost:10080/v1/keyval/mykey -H "Content-Type: text/plain" -d '{"value": "myval1"}'

# Read a Key
curl -vvv -X GET -k https://localhost:10080/v1/keyval/mykey -H "Content-Type: text/plain"

# Update a Key
curl -vvv -X POST -k https://localhost:10080/v1/keyval/mykey -H "Content-Type: text/plain" -d '{"value": "myval2"}'

# Delete a Key
curl -vvv -X DELETE -k https://localhost:10080/v1/keyval/mykey -H "Content-Type: text/plain"

### View the Swagger UI

# View the swagger-ui
curl -k https://localhost:10080/swagger-ui/

# View the swagger file that has been generated from GRPC
curl -k https://localhost:10080/swagger.json
```
