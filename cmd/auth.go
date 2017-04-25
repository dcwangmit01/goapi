package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"

	jwt "github.com/dcwangmit01/goapi/app/jwt"
	pb "github.com/dcwangmit01/goapi/app/pb"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Obtain JWT token and print it out",
	Long: `
Auth:
    goapi auth <email> <password>
`,
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRunAuth(authAndPrint)
	},
}

func authAndPrint(client pb.AuthClient, ctx context.Context) {
	req := &pb.AuthRequestMessage{
		"password",
		os.Args[2],
		os.Args[3]}

	logutil.AddCtx(log.WithFields(log.Fields{
		"message": req,
	})).Info("Sent RPC Request")

	rsp, _ := client.Auth(ctx, req)

	logutil.AddCtx(log.WithFields(log.Fields{
		"message": rsp,
	})).Info("Received RPC Reply")

	tokenStr := rsp.GetAccessToken()
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
