package service

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var Registry = ServiceRegistry{}

type GrpcServiceHandler func(*grpc.Server)
type GrpcGatewayHandler func(context.Context,
	*runtime.ServeMux, string, []grpc.DialOption) error

type ServiceRegistry struct {
	GrpcServiceHandlers []GrpcServiceHandler
	GrpcGatewayHandlers []GrpcGatewayHandler
}

func (r *ServiceRegistry) AddGrpcServiceHandler(f GrpcServiceHandler) {
	r.GrpcServiceHandlers = append(r.GrpcServiceHandlers, f)
}

func (r *ServiceRegistry) AddGrpcGatewayHandler(f GrpcGatewayHandler) {
	r.GrpcGatewayHandlers = append(r.GrpcGatewayHandlers, f)
}
