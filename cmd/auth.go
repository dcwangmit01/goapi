package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/grpc-gw-poc/app/logutil"

	pb "github.com/dcwangmit01/grpc-gw-poc/app"
	jwt "github.com/dcwangmit01/grpc-gw-poc/app/jwt"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Obtain JWT token and print it out",
	Long: `
Auth:
    grpc-gw-poc auth <email> <password>
`,
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRun(authAndPrint)
	},
}

func authAndPrint(client pb.AppClient, ctx context.Context) {
	req := &pb.AuthRequestMessage{os.Args[2], os.Args[3]}

	logutil.AddCtx(log.WithFields(log.Fields{
		"message": req,
	})).Info("Sent RPC Request")

	rsp, _ := client.Auth(ctx, req)

	logutil.AddCtx(log.WithFields(log.Fields{
		"message": rsp,
	})).Info("Received RPC Reply")

	tokenStr := rsp.GetToken()
	token, customClaims, err := jwt.ParseJwt(tokenStr)

	logutil.AddCtx(log.WithFields(log.Fields{
		"valid":  token.Valid,
		"claims": customClaims,
		"error":  err,
	})).Info("Parsed JWT")
}

func init() {
	RootCmd.AddCommand(authCmd)
}
