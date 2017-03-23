package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/dcwangmit01/grpc-gw-poc/app"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Example client gRPC service CLI client",
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		creds := credentials.NewClientTLSFromCert(certPool, "localhost:10080")
		opts = append(opts, grpc.WithTransportCredentials(creds))
		conn, err := grpc.Dial(serverAddress, opts...)
		if err != nil {
			grpclog.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := pb.NewAppClient(conn)

		msg, err := client.Echo(context.Background(), &pb.EchoMessage{strings.Join(os.Args[2:], " ")})
		println(msg.Value)

	},
}

func init() {
	RootCmd.AddCommand(clientCmd)
}

