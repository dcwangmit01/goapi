package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/grpc-gw-poc/app/logutil"

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
}

var keyvalCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRun(keyvalCreate)
	},
}

func keyvalCreate(client pb.AppClient) {
	msg, _ := client.KeyValCreate(context.Background(), &pb.KeyValMessage{os.Args[3], os.Args[4]})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

var keyvalReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRun(keyvalRead)
	},
}

func keyvalRead(client pb.AppClient) {
	msg, _ := client.KeyValRead(context.Background(), &pb.KeyValMessage{os.Args[3], ""})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

var keyvalUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRun(keyvalUpdate)
	},
}

func keyvalUpdate(client pb.AppClient) {
	msg, _ := client.KeyValUpdate(context.Background(), &pb.KeyValMessage{os.Args[3], os.Args[4]})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

var keyvalDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRun(keyvalDelete)
	},
}

func keyvalDelete(client pb.AppClient) {
	msg, _ := client.KeyValDelete(context.Background(), &pb.KeyValMessage{os.Args[3], ""})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

func init() {
	RootCmd.AddCommand(keyvalRootCmd)
	keyvalRootCmd.AddCommand(keyvalCreateCmd)
	keyvalRootCmd.AddCommand(keyvalReadCmd)
	keyvalRootCmd.AddCommand(keyvalUpdateCmd)
	keyvalRootCmd.AddCommand(keyvalDeleteCmd)
}
