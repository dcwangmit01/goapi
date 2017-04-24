# goapi

This project is a means to learn golang and play with GRPC, JWT, etc.

This Proof of Concept provides code for:
* Generating resources that were included in grpc-gateway-example, without build scripts
  * GRPC files including protobuf, grpc, gateway, swagger
  * SSL certs using [cfssl](https://github.com/cloudflare/cfssl)
  * Swagger go-bindata-assetfs files for swagger-ui and the swagger file
* Using a vendored golang package management framework (glide)
* Writing a golang CLI tool using [cobra](https://github.com/spf13/cobra)
* Cross-compiling statically linked Linux binaries for amd64 and armv7.

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
