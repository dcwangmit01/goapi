package service

import (
	"errors"

	context "golang.org/x/net/context"
	metadata "google.golang.org/grpc/metadata"

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

	// find the user
	u, err := ac.GetUserByEmail(in.GetEmail())
	if err != nil {
		logutil.AddCtx(log.WithFields(log.Fields{
			"message": in,
			"error":   err,
		})).Warn("Auth RPC Failed")
		return &pb.AuthResponseMessage{}, err
	}

	// validate the user password
	err = u.ValidatePassword(in.GetPassword())
	if err != nil {
		logutil.AddCtx(log.WithFields(log.Fields{
			"message": in,
			"error":   err,
		})).Warn("Auth RPC Failed")
		return &pb.AuthResponseMessage{}, err
	}

	// create the JWT token, which contains claims
	jwtStr, err := u.GenerateJwt()
	if err != nil {
		logutil.AddCtx(log.WithFields(log.Fields{
			"message": in,
			"error":   err,
		})).Warn("Auth RPC Failed")
		return &pb.AuthResponseMessage{}, err
	}

	return &pb.AuthResponseMessage{Token: jwtStr}, nil
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
	md, ok := metadata.FromContext(c)
	if !ok {
		return m, errors.New("Cannot decode metadata")
	}
	logutil.AddCtx(log.WithFields(log.Fields{
		"message":  m,
		"metadata": md,
	})).Info("Received RPC Request")
	if !kv.SqlKV.HasKey(m.Key) {
		return m, errors.New("Cannot read non-existent Key")
	}
	m.Value = kv.SqlKV.String(m.Key)
	return m, nil
}

func (s *myService) KeyValUpdate(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")
	if !kv.SqlKV.HasKey(m.Key) {
		return &pb.EmptyMessage{}, errors.New("Cannot update non-existent Key")
	}
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
