package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"

	pb "github.com/dcwangmit01/goapi/app"
	config "github.com/dcwangmit01/goapi/app/config"
	svc "github.com/dcwangmit01/goapi/app/service"
	kv "github.com/dcwangmit01/goapi/app/sqlitekv"
	certs "github.com/dcwangmit01/goapi/resources/certs"
	swf "github.com/dcwangmit01/goapi/resources/swagger/files"
	sw "github.com/dcwangmit01/goapi/resources/swagger/ui"
)

/* Overview

Request comes in on http.Server.  http.Server calls grpcHandlerFunc which
determines via http protocol and content-type header whether the request is a
GRPC request, or an http request.

If the request is a GRPC request, then incoming request is passed to the
grpcServer handler to service the request.  If it is not a GRPC request and
instead a web/rest request, then the incoming request is passed to the top
level "mux" multiplexor handler.

The top level "mux" multiplexor handler handles HTTP requests.  URI paths are
matched, and requests are passed to the corresponding handlers.  /swagger.json
is served by a handler function, /swagger-ui is served by a http.Fileserver
handler, and all other URI paths "/" are passed to the grpc-gw mux.  The
grpc-gw mux is a handler that matches incoming JSON REST requests, and proxies
them to the grpcServer handler.

*/

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks
		// https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61

		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    sw.Asset,
		AssetDir: sw.AssetDir,
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func StartServer() {

	kv.Init() // initialize the sql key/value database

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certs.CertPool, config.ServerAddress))}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServer(grpcServer, svc.NewAuthService())     // grpc
	pb.RegisterKeyValServer(grpcServer, svc.NewKeyValService()) // grpc
	ctx := context.Background()

	// client credentials
	ccreds := credentials.NewTLS(&tls.Config{
		ServerName: config.ServerAddress,
		RootCAs:    certs.CertPool,
	})

	// client options
	copts := []grpc.DialOption{grpc.WithTransportCredentials(ccreds)}

	data, _ := swf.Asset("swagger.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(data))
	})

	gwmux := runtime.NewServeMux()
	// Registers function handlers for each uri pattern defined by grpc-gw
	//   matches http requests to patterns and invokes the corresponding handler
	var err error
	err = pb.RegisterAuthHandlerFromEndpoint(ctx, gwmux, config.ServerAddress, copts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}
	err = pb.RegisterKeyValHandlerFromEndpoint(ctx, gwmux, config.ServerAddress, copts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	mux.Handle("/", gwmux)
	serveSwagger(mux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*certs.KeyPair},
			NextProtos:   []string{"h2"},
		},
	}

	fmt.Printf("grpc on port: %d\n", config.Port)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		logutil.AddCtx(log.WithFields(log.Fields{
			"error": err,
		})).Info("ListenAndServe")
	}

	return
}
