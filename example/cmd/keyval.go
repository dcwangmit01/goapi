package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dcwangmit01/goapi/example/config"
	"github.com/dcwangmit01/goapi/resources/certs"
	"github.com/dcwangmit01/goapi/util"

	clt "github.com/dcwangmit01/goapi/client"
	pb "github.com/dcwangmit01/goapi/example/pb"
)

var (
	optionUsername string
	optionPassword string
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
	Example: `  # Create a key:
  goapi keyval create <key> <value>

  # Read a key:
  goapi keyval read <key>

  # Update a key:
  goapi keyval update <key> <value>

  # Delete a key:
  goapi keyval delete <key>`,
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

	// authenticate
	tokenStr, err := clt.Authenticate(optionUsername, optionPassword)
	if err != nil {
		return err
	}

	// connect with the jwt auth token
	conn, ctx, err := clt.ConnectWithToken(config.Host, config.Port, tokenStr, certs.CertPool)
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
	client := pb.NewKeyValClient(conn)

	var rsp interface{}
	switch operation {
	case kvCreate:
		rsp, err = client.KeyValCreate(ctx, req)
	case kvRead:
		rsp, err = client.KeyValRead(ctx, req)
	case kvUpdate:
		rsp, err = client.KeyValUpdate(ctx, req)
	case kvDelete:
		rsp, err = client.KeyValDelete(ctx, req)
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
