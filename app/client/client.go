package client

import (
	"crypto/x509"
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"

	"github.com/dcwangmit01/goapi/app/config"
	pb "github.com/dcwangmit01/goapi/app/pb"
	"github.com/dcwangmit01/goapi/resources/certs"
)

// returns a connection object with established connection, if successful
func GrpcTlsConnect(host string, port int, ctx context.Context, certPool *x509.CertPool) (*grpc.ClientConn, error) {

	serverAddress := fmt.Sprintf("%s:%d", host, port)

	// Create the TLS connection optiosn
	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certs.CertPool, serverAddress)
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// Dial (aka. Connect) to the GRPC server
	conn, err := grpc.Dial(serverAddress, opts...)

	// return the connection and error directly, relying on the
	// caller to close the connection
	return conn, err
}

// returns a JWT auth token, if successful
func Authenticate(username string, password string) (string, error) {

	// connect to the grpc server
	ctx := context.Background()
	conn, err := GrpcTlsConnect(config.Host, config.Port, ctx, certs.CertPool)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// construct the request
	req := &pb.AuthRequestMessage{
		GrantType: "password", // must be 'password'
		Username:  username,
		Password:  password,
	}

	// create the client and send the request
	client := pb.NewAuthClient(conn)
	rsp, err := client.Auth(ctx, req)
	if err != nil {
		return "", err
	}

	// log the info
	logutil.AddCtx(log.WithFields(log.Fields{
		"request":  req,
		"response": rsp,
	})).Debug("RPC Request and Response")

	// validate the token string and return it
	tokenString := rsp.GetAccessToken()
	if tokenString == "" {
		return "", errors.New("Recevied empty token string")
	}
	return rsp.GetAccessToken(), nil
}

func StructToYamlStr(s interface{}) (string, error) {

	dump, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(dump), err
}
