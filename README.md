# grpc-gw-poc

This project is a means to learn golang and play with GRPC.  This
application code is based off of
[grpc-gateway-example](https://github.com/philips/grpc-gateway-example).
The build scripts and framework are unique to this project.

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
./bin/linux_amd64/grpc-gw-poc

# Start the GRPC+JSON server on port 10080
./bin/linux_amd64/grpc-gw-poc serve

# Run the GRPC client
./bin/linux_amd64/grpc-gw-poc client "foo bar"

# Run the curl client against the gateway
curl -vvv -X POST -k https://localhost:10080/v1/echo -H "Content-Type: text/plain" -d '{"value": "foo bar"}'

# View the swagger-ui
curl -k https://localhost:10080/swagger-ui/

# View the swagger file that has been generated from GRPC
curl -k https://localhost:10080/swagger.json
```
