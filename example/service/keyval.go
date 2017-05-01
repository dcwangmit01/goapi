package service

import (
	"errors"
	"fmt"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/dcwangmit01/goapi/example/config"
	pb "github.com/dcwangmit01/goapi/example/pb"
	"github.com/dcwangmit01/goapi/registry"
	"github.com/dcwangmit01/goapi/sqlitekv"
	"github.com/dcwangmit01/goapi/util"
)

func init() {
	registry.ServiceRegistry.AddGrpcGatewayHandler(pb.RegisterKeyValHandlerFromEndpoint)
	registry.ServiceRegistry.AddGrpcServiceHandler(func(grpcServer *grpc.Server) {
		pb.RegisterKeyValServer(grpcServer, NewKeyValService())
	})
}

var skv = sqlitekv.New(config.AppName + ".db")

type kvService struct{}

func (s *kvService) KeyValCreate(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if skv.HasKey(m.Key) {
		return &pb.KeyValMessage{}, errors.New("Cannot create existing Key")
	}
	skv.SetString(m.Key, m.Value)
	return m, nil
}

func (s *kvService) KeyValRead(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {

	// Example code on accessing user
	/////////////////////////////////////////////////////////////////////
	user, ok := config.UserFromContext(c)
	if !ok {
		return m, errors.New("Unable to locate user")
	}
	dump, _ := util.StructToYamlStr(user)
	fmt.Println(dump)
	/////////////////////////////////////////////////////////////////////

	if !skv.HasKey(m.Key) {
		return m, errors.New("Cannot read non-existent Key")
	}
	m.Value = skv.String(m.Key)
	return m, nil
}

func (s *kvService) KeyValUpdate(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if !skv.HasKey(m.Key) {
		return &pb.KeyValMessage{}, errors.New("Cannot update non-existent Key")
	}
	skv.SetString(m.Key, m.Value)
	return m, nil
}

func (s *kvService) KeyValDelete(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if !skv.HasKey(m.Key) {
		return &pb.KeyValMessage{}, errors.New("Cannot delete non-existent Key")
	}
	m.Value = skv.String(m.Key)
	skv.Del(m.Key)
	return m, nil
}

func NewKeyValService() *kvService {
	return new(kvService)
}
