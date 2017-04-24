package service

import (
	"errors"

	context "golang.org/x/net/context"
	metadata "google.golang.org/grpc/metadata"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"

	pb "github.com/dcwangmit01/goapi/app"
	kv "github.com/dcwangmit01/goapi/app/sqlitekv"
)

type kvService struct{}

func (s *kvService) KeyValCreate(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")

	if kv.SqlKV.HasKey(m.Key) {
		return &pb.EmptyMessage{}, errors.New("Cannot create existing Key")
	}
	kv.SqlKV.SetString(m.Key, m.Value)
	return &pb.EmptyMessage{}, nil
}

func (s *kvService) KeyValRead(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
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

func (s *kvService) KeyValUpdate(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")
	if !kv.SqlKV.HasKey(m.Key) {
		return &pb.EmptyMessage{}, errors.New("Cannot update non-existent Key")
	}
	kv.SqlKV.SetString(m.Key, m.Value)
	return &pb.EmptyMessage{}, nil
}

func (s *kvService) KeyValDelete(c context.Context, m *pb.KeyValMessage) (*pb.EmptyMessage, error) {
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": m,
	})).Info("Received RPC Request")
	if !kv.SqlKV.HasKey(m.Key) {
		return &pb.EmptyMessage{}, errors.New("Cannot delete non-existent Key")
	}
	kv.SqlKV.Del(m.Key)
	return &pb.EmptyMessage{}, nil
}

func NewKeyValService() *kvService {
	return new(kvService)
}
