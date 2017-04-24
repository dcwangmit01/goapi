package service

import (
	context "golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"

	pb "github.com/dcwangmit01/goapi/app"
	cnf "github.com/dcwangmit01/goapi/app/config"
)

type authService struct{}

func (s *authService) Auth(ctx context.Context, in *pb.AuthRequestMessage) (*pb.AuthResponseMessage, error) {
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

func NewAuthService() *authService {
	return new(authService)
}
