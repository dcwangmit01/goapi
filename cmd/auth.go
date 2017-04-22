package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/grpc-gw-poc/app/logutil"

	pb "github.com/dcwangmit01/grpc-gw-poc/app"
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

func authAndPrint(client pb.AppClient) {
	req := &pb.AuthRequestMessage{os.Args[2], os.Args[3]}

	logutil.AddCtx(log.WithFields(log.Fields{
		"message": req,
	})).Info("Sent RPC Request")

	rsp, _ := client.Auth(context.Background(), req)

	logutil.AddCtx(log.WithFields(log.Fields{
		"message": rsp,
	})).Info("Received RPC Reply")
}

func init() {
	RootCmd.AddCommand(authCmd)
}
