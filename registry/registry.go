package registry

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var ServiceRegistry = serviceRegistry{}

type GrpcServiceHandler func(*grpc.Server)
type GrpcGatewayHandler func(context.Context,
	*runtime.ServeMux, string, []grpc.DialOption) error

type serviceRegistry struct {
	GrpcServiceHandlers []GrpcServiceHandler
	GrpcGatewayHandlers []GrpcGatewayHandler
}

func (sr *serviceRegistry) AddGrpcServiceHandler(f GrpcServiceHandler) {
	sr.GrpcServiceHandlers = append(sr.GrpcServiceHandlers, f)
}

func (sr *serviceRegistry) AddGrpcGatewayHandler(f GrpcGatewayHandler) {
	sr.GrpcGatewayHandlers = append(sr.GrpcGatewayHandlers, f)
}
