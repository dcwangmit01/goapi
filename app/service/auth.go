package service

import (
	"errors"

	context "golang.org/x/net/context"

	cnf "github.com/dcwangmit01/goapi/app/config"
	pb "github.com/dcwangmit01/goapi/app/pb"
)

type authService struct{}

func (s *authService) Auth(ctx context.Context, in *pb.AuthRequestMessage) (*pb.AuthResponseMessage, error) {

	ac := cnf.SingletonAppConfig

	// check the grant type is "password".
	//   required per oauth2 spec Client Credentials Grant Type
	if in.GetGrantType() != "password" {
		return &pb.AuthResponseMessage{}, errors.New("Grant type must be 'password'")
	}

	// find the user
	u, err := ac.GetUserByUsername(in.GetUsername())
	if err != nil {
		return &pb.AuthResponseMessage{}, err
	}

	// validate the user password
	err = u.ValidatePassword(in.GetPassword())
	if err != nil {
		return &pb.AuthResponseMessage{}, err
	}

	// create the JWT token, which contains claims
	duration := int64(3600) // 1 hour
	jwtStr, err := u.GenerateJwt(duration)
	if err != nil {
		return &pb.AuthResponseMessage{}, err
	}

	return &pb.AuthResponseMessage{
		AccessToken: jwtStr,
		TokenType:   "JWT",
		ExpiresIn:   duration,
	}, nil
}

func NewAuthService() *authService {
	return new(authService)
}
