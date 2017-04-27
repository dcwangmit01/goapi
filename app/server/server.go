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

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_gw_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	config "github.com/dcwangmit01/goapi/app/config"
	pb "github.com/dcwangmit01/goapi/app/pb"
	svc "github.com/dcwangmit01/goapi/app/service"
	kv "github.com/dcwangmit01/goapi/app/sqlitekv"
	certs "github.com/dcwangmit01/goapi/resources/certs"
	swf "github.com/dcwangmit01/goapi/resources/swagger/files"
	sw "github.com/dcwangmit01/goapi/resources/swagger/ui"
)

/* Overview

Request comes in on http.Server.  http.Server calls triageHandlerFunc which
determines via http protocol and content-type header whether the request is a
GRPC request, or an http request.

If the request is a GRPC request, then incoming request is passed to
the grpcServer handler to service the request.  If it is not a GRPC
request and instead a web/rest request, then the incoming request is
passed to the webHandler (top level "mux" multiplexor)

The top level "mux" multiplexor handler handles HTTP requests.  URI paths are
matched, and requests are passed to the corresponding handlers.  /swagger.json
is served by a handler function, /swagger-ui is served by a http.Fileserver
handler, and all other URI paths "/" are passed to the grpc-gw mux.  The
grpc-gw mux is a handler that matches incoming JSON REST requests, and proxies
them to the grpcServer handler.

Interesting reads:
http://www.alexedwards.net/blog/a-recap-of-request-handling
http://www.alexedwards.net/blog/making-and-using-middleware
*/

func triageHandlerFunc(grpcHandler http.Handler, webHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcHandler.ServeHTTP(w, r)
		} else {
			webHandler.ServeHTTP(w, r)
		}
	})
}

func registerSwaggerFileHandler(mux *http.ServeMux) {
	data, _ := swf.Asset("swagger.json")
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(data))
	})
}

func registerSwaggerUiHandler(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose swagger-ui files on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    sw.Asset,
		AssetDir: sw.AssetDir,
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func registerGrpcGatewayHandlers(mux *http.ServeMux) {
	var err error

	gwmux := grpc_gw_runtime.NewServeMux()
	ctx := context.Background()
	ccreds := credentials.NewTLS(&tls.Config{
		ServerName: config.ServerAddress,
		RootCAs:    certs.CertPool,
	})
	copts := []grpc.DialOption{grpc.WithTransportCredentials(ccreds)}
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
}

func StartServer() {

	kv.Init() // initialize the sql key/value database

	/*
	   Create the grpc handler
	*/
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certs.CertPool, config.ServerAddress)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logger))),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServer(grpcServer, svc.NewAuthService())
	pb.RegisterKeyValServer(grpcServer, svc.NewKeyValService())

	/*
	   Create the web handler
	*/
	mux := http.NewServeMux()
	registerSwaggerFileHandler(mux)
	registerSwaggerUiHandler(mux)
	registerGrpcGatewayHandlers(mux)

	/*
	   Start the server
	*/
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: triageHandlerFunc(CommonMiddleware.Then(grpcServer), CommonMiddleware.Then(mux)),
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
