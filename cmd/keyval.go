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
		grpcDialAndRunKeyVal(keyvalCreate)
	},
}

func keyvalCreate(client pb.KeyValClient, ctx context.Context) {
	msg, _ := client.KeyValCreate(ctx, &pb.KeyValMessage{os.Args[3], os.Args[4]})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

var keyvalReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRunKeyVal(keyvalRead)
	},
}

func keyvalRead(client pb.KeyValClient, ctx context.Context) {
	msg, _ := client.KeyValRead(ctx, &pb.KeyValMessage{os.Args[3], ""})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

var keyvalUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRunKeyVal(keyvalUpdate)
	},
}

func keyvalUpdate(client pb.KeyValClient, ctx context.Context) {
	msg, _ := client.KeyValUpdate(ctx, &pb.KeyValMessage{os.Args[3], os.Args[4]})
	logutil.AddCtx(log.WithFields(log.Fields{
		"message": msg,
	})).Info("Sent RPC Request")
}

var keyvalDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete new Key/Value on gRPC service",
	Run: func(cmd *cobra.Command, args []string) {
		grpcDialAndRunKeyVal(keyvalDelete)
	},
}

func keyvalDelete(client pb.KeyValClient, ctx context.Context) {
	msg, _ := client.KeyValDelete(ctx, &pb.KeyValMessage{os.Args[3], ""})
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
