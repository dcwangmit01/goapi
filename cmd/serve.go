package cmd

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"strings"
	"bytes"

	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/elazarl/go-bindata-assetfs"

	pb "github.com/dcwangmit01/grpc-gw-poc/app"
	sw "github.com/dcwangmit01/grpc-gw-poc/resources/swagger/ui"
	swf "github.com/dcwangmit01/grpc-gw-poc/resources/swagger/files"
	certs "github.com/dcwangmit01/grpc-gw-poc/resources/certs"
)

var (
	keyPair *tls.Certificate
	certPool *x509.CertPool
	serverAddress string
	host = "localhost"
	port = 10080
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches the example webserver on https://localhost:10000",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

type myService struct{}

func (m *myService) Echo(c context.Context, s *pb.EchoMessage) (*pb.EchoMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

func newServer() *myService {
	return new(myService)
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
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


func init() {

	var err error

	key, err := certs.Asset("insecure-key.pem")
	if err != nil {
		panic(err)
	}

	pem, err := certs.Asset("insecure.pem")
	if err != nil {
		panic(err)
	}

	pair, err := tls.X509KeyPair(pem, key)
	if err != nil {
		panic(err)
	}
	keyPair = &pair

        certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(pem)
	if !ok {
		panic("bad certs")
	}

	serverAddress = fmt.Sprintf("%s:%d", host, port)
}

func serve() {

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, serverAddress))}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAppServer(grpcServer, newServer())
	ctx := context.Background()

	// client credentials
	ccreds := credentials.NewTLS(&tls.Config{
		ServerName: serverAddress,
		RootCAs:    certPool,
	})

	// client options
	copts := []grpc.DialOption{grpc.WithTransportCredentials(ccreds)}

	data, _ := swf.Asset("swagger.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(data))
	})

	gwmux := runtime.NewServeMux()
	err := pb.RegisterAppHandlerFromEndpoint(ctx, gwmux, serverAddress, copts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	mux.Handle("/", gwmux)
	serveSwagger(mux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    serverAddress,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	fmt.Printf("grpc on port: %d\n", port)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}
