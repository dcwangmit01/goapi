package service

import (
	"errors"

	context "golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/grpc-gw-poc/app/logutil"

	pb "github.com/dcwangmit01/grpc-gw-poc/app"
	cnf "github.com/dcwangmit01/grpc-gw-poc/app/config"
	kv "github.com/dcwangmit01/grpc-gw-poc/app/sqlitekv"
)

type myService struct{}

func (s *myService) Auth(ctx context.Context, in *pb.AuthRequestMessage) (*pb.AuthResponseMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": in,
	})).Info("Received RPC Request")

	ac := cnf.SingletonAppConfig
	u := ac.GetUserByEmail(in.GetEmail())
	if u == nil {
		return &pb.AuthResponseMessage{}, errors.New("User Not Found")
	}

	err := u.ValidatePassword(in.GetPassword())
	if err != nil {
		return &pb.AuthResponseMessage{}, errors.New("Invalid Password")
	}

	return &pb.AuthResponseMessage{Token: "a new JWT token"}, nil
}

func (s *myService) KeyValCreate(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")

	if kv.SqlKV.HasKey(m.Key) {
		return &pb.EmptyMessage{}, errors.New("Cannot create existing Key")
	}

	kv.SqlKV.SetString(m.Key, m.Value)
	return &pb.EmptyMessage{}, nil
}

func (s *myService) KeyValRead(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")
	m.Value = kv.SqlKV.String(m.Key)
	return m, nil
}

func (s *myService) KeyValUpdate(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")
	kv.SqlKV.SetString(m.Key, m.Value)
	return &pb.EmptyMessage{}, nil
}

func (s *myService) KeyValDelete(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")
	if !kv.SqlKV.HasKey(m.Key) {
		return &pb.EmptyMessage{}, errors.New("Cannot delete non-existent Key")
	}
	kv.SqlKV.Del(m.Key)
	return &pb.EmptyMessage{}, nil
}

func NewServer() *myService {
	return new(myService)
}
