package service

import (
	"errors"
	"fmt"

	context "golang.org/x/net/context"

	"github.com/dcwangmit01/goapi/app/client"
	"github.com/dcwangmit01/goapi/app/config"
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

	// Example code on accessing user
	/////////////////////////////////////////////////////////////////////
	user, ok := config.UserFromContext(c)
	if !ok {
		return m, errors.New("Unable to locate user")
	}
	dump, _ := client.StructToYamlStr(user)
	fmt.Println(dump)
	/////////////////////////////////////////////////////////////////////

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
