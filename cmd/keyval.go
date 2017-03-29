package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/dcwangmit01/grpc-gw-poc/app"
)

var keyvalRootCmd = &cobra.Command{
	Use:   "keyval",
	Short: "Client used to set Key/Value on gRPC service",
	Long: `
Create a key;
    grpc-gw-poc keyval create <key> <value>

Read a key:
    grpc-gw-poc keyval read <key>

Update a key:
    grpc-gw-poc keyval update <key> <value>

Delete a key:
    grpc-gw-poc keyval delete <key>
`,
};

var keyvalCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		connectClientRunFunc(keyvalCreate)
	},
};

func connectClientRunFunc(f func(pb.AppClient)) {
	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certPool, "localhost:10080")
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAppClient(conn)
	f(client)
}

func keyvalCreate(client pb.AppClient) {
	msg, _ := client.KeyValCreate(context.Background(), &pb.KeyValMessage{os.Args[3], os.Args[4]})
	fmt.Printf("rpc client request s(%q)\n", msg)
}

func init() {
	RootCmd.AddCommand(keyvalRootCmd)
	keyvalRootCmd.AddCommand(keyvalCreateCmd)
}

