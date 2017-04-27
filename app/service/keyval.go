package service

import (
	"errors"

	context "golang.org/x/net/context"

	pb "github.com/dcwangmit01/goapi/app/pb"
	kv "github.com/dcwangmit01/goapi/app/sqlitekv"
)

type kvService struct{}

func (s *kvService) KeyValCreate(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if kv.SqlKV.HasKey(m.Key) {
		return &pb.KeyValMessage{}, errors.New("Cannot create existing Key")
	}
	kv.SqlKV.SetString(m.Key, m.Value)
	return m, nil
}

func (s *kvService) KeyValRead(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if !kv.SqlKV.HasKey(m.Key) {
		return m, errors.New("Cannot read non-existent Key")
	}
	m.Value = kv.SqlKV.String(m.Key)
	return m, nil
}

func (s *kvService) KeyValUpdate(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if !kv.SqlKV.HasKey(m.Key) {
		return &pb.KeyValMessage{}, errors.New("Cannot update non-existent Key")
	}
	kv.SqlKV.SetString(m.Key, m.Value)
	return m, nil
}

func (s *kvService) KeyValDelete(c context.Context, m *pb.KeyValMessage) (*pb.KeyValMessage, error) {
	if !kv.SqlKV.HasKey(m.Key) {
		return &pb.KeyValMessage{}, errors.New("Cannot delete non-existent Key")
	}
	m.Value = kv.SqlKV.String(m.Key)
	kv.SqlKV.Del(m.Key)
	return m, nil
}

func NewKeyValService() *kvService {
	return new(kvService)
}
