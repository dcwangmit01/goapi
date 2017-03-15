package main

import (
 "flag"
 "net/http"

 "github.com/golang/glog"
 "golang.org/x/net/context"
 "github.com/grpc-ecosystem/grpc-gateway/runtime"
 "google.golang.org/grpc"

// gw "path/to/your_service_package"
  gw "github.com/dcwangmit01/pax-api/entry-lib"
)

var (
 echoEndpoint = flag.String("echo_endpoint", "localhost:9090", "endpoint of YourService")
)

func run() error {
 ctx := context.Background()
 ctx, cancel := context.WithCancel(ctx)
 defer cancel()

 mux := runtime.NewServeMux()
 opts := []grpc.DialOption{grpc.WithInsecure()}
 err := gw.RegisterPaxEntryApiHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
 if err != nil {
   return err
 }

 return http.ListenAndServe(":8080", mux)
}

func main() {
 flag.Parse()
 defer glog.Flush()

 if err := run(); err != nil {
   glog.Fatal(err)
 }
}
