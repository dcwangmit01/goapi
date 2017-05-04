package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dcwangmit01/goapi/config"
	"github.com/dcwangmit01/goapi/resources/certs"
	"github.com/dcwangmit01/goapi/util"

	"github.com/dcwangmit01/goapi/client"
	pb "github.com/dcwangmit01/goapi/example/pb"
)

func init() {
	RootCmd.AddCommand(keyvalRootCmd)
	keyvalRootCmd.AddCommand(keyvalCreateCmd)
	keyvalRootCmd.AddCommand(keyvalReadCmd)
	keyvalRootCmd.AddCommand(keyvalUpdateCmd)
	keyvalRootCmd.AddCommand(keyvalDeleteCmd)
	keyvalRootCmd.PersistentFlags().StringVarP(&optionUsername, "username", "u", "", "Username for authentication")
	keyvalRootCmd.PersistentFlags().StringVarP(&optionPassword, "password", "p", "", "Password for authentication")
}

type RestOp int

const (
	kvCreate RestOp = iota
	kvRead   RestOp = iota
	kvUpdate RestOp = iota
	kvDelete RestOp = iota
)

var keyvalRootCmd = &cobra.Command{
	Use:   "keyval",
	Short: "Client used to set Key/Value on gRPC service",
	Example: fmt.Sprintf(`  # Create a key:
  %[1]v keyval create <key> <value>

  # Read a key:
  %[1]v keyval read <key>

  # Update a key:
  %[1]v keyval update <key> <value>

  # Delete a key:
  %[1]v keyval delete <key>`, config.GetAppName()),
}

var keyvalCreateCmd = &cobra.Command{
	Use:   "create <key> <value>",
	Short: "Create new Key/Value on gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return keyvalCreate(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

var keyvalReadCmd = &cobra.Command{
	Use:   "read <key>",
	Short: "Read new Key/Value on gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return keyvalRead(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

var keyvalUpdateCmd = &cobra.Command{
	Use:   "update <key> <value>",
	Short: "Update new Key/Value on gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return keyvalUpdate(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

var keyvalDeleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete new Key/Value on gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return keyvalDelete(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

func keyvalCreate(cmd *cobra.Command, args []string) error {
	return keyvalHelper(cmd, args, 2, 0, 1, kvCreate)
}

func keyvalRead(cmd *cobra.Command, args []string) error {
	return keyvalHelper(cmd, args, 1, 0, -1, kvRead)
}

func keyvalUpdate(cmd *cobra.Command, args []string) error {
	return keyvalHelper(cmd, args, 2, 0, 1, kvUpdate)
}

func keyvalDelete(cmd *cobra.Command, args []string) error {
	return keyvalHelper(cmd, args, 1, 0, -1, kvDelete)
}

func keyvalHelper(cmd *cobra.Command, args []string,
	expectedNumArgs int, keyArgIndex int, valueArgIndex int, operation RestOp) error {

	// validate args
	if len(args) != expectedNumArgs {
		cmd.Usage()
		return invalidInputErr
	}

	// get the auth token
	tokenStr, err := client.GetAuthTokenFromOptionOrConfigOrStdin(
		true,           // optionTry
		optionUsername, // optionUsername
		optionPassword, // optionPassword
		true,           // optionContinue
		true,           // configTry
		true,           // configContinue
		false,          // stdinTry
		false,          // stdinContinue
		true,           // saveNewToken
	)
	if err != nil {
		return err
	}

	// connect with the jwt auth token
	conn, ctx, err := client.ConnectWithToken(config.GetHost(), config.GetPort(), tokenStr, certs.CertPool)
	if err != nil {
		return err
	}
	defer conn.Close()

	// construct the request
	key := ""
	if keyArgIndex >= 0 {
		key = args[keyArgIndex]
	}
	value := ""
	if valueArgIndex >= 0 {
		value = args[valueArgIndex]
	}
	req := &pb.KeyValMessage{
		Key:   key,
		Value: value,
	}

	// create the client and send the request
	cli := pb.NewKeyValClient(conn)

	var rsp interface{}
	switch operation {
	case kvCreate:
		rsp, err = cli.KeyValCreate(ctx, req)
	case kvRead:
		rsp, err = cli.KeyValRead(ctx, req)
	case kvUpdate:
		rsp, err = cli.KeyValUpdate(ctx, req)
	case kvDelete:
		rsp, err = cli.KeyValDelete(ctx, req)
	default:
		panic("Code Error")
	}
	if err != nil {
		return err
	}

	// print the response to stdout
	dump, err := util.StructToYamlStr(rsp)
	if err != nil {
		return err
	}
	fmt.Printf("%v", dump)

	return nil
}
