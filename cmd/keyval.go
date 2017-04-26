package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dcwangmit01/goapi/app/config"
	"github.com/dcwangmit01/goapi/resources/certs"

	clt "github.com/dcwangmit01/goapi/app/client"
	pb "github.com/dcwangmit01/goapi/app/pb"
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

	// validate args
	if len(args) != 2 {
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
	req := &pb.KeyValMessage{
		Key:   args[0],
		Value: args[1],
	}

	// create the client and send the request
	client := pb.NewKeyValClient(conn)
	rsp, err := client.KeyValCreate(ctx, req)

	// print the response to stdout
	dump, err := clt.StructToYamlStr(rsp)
	if err != nil {
		return err
	}
	fmt.Printf("%v", dump)

	return nil
}

func keyvalRead(cmd *cobra.Command, args []string) error {

	// validate args
	if len(args) != 1 {
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
	req := &pb.KeyValMessage{
		Key:   args[0],
		Value: "",
	}

	// create the client and send the request
	client := pb.NewKeyValClient(conn)
	rsp, err := client.KeyValRead(ctx, req)

	// print the response to stdout
	dump, err := clt.StructToYamlStr(rsp)
	if err != nil {
		return err
	}
	fmt.Printf("%v", dump)

	return nil
}

func keyvalUpdate(cmd *cobra.Command, args []string) error {

	// validate args
	if len(args) != 2 {
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
	req := &pb.KeyValMessage{
		Key:   args[0],
		Value: args[1],
	}

	// create the client and send the request
	client := pb.NewKeyValClient(conn)
	rsp, err := client.KeyValUpdate(ctx, req)

	// print the response to stdout
	dump, err := clt.StructToYamlStr(rsp)
	if err != nil {
		return err
	}
	fmt.Printf("%v", dump)

	return nil
}

func keyvalDelete(cmd *cobra.Command, args []string) error {

	// validate args
	if len(args) != 1 {
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
	req := &pb.KeyValMessage{
		Key:   args[0],
		Value: "",
	}

	// create the client and send the request
	client := pb.NewKeyValClient(conn)
	rsp, err := client.KeyValDelete(ctx, req)

	// print the response to stdout
	dump, err := clt.StructToYamlStr(rsp)
	if err != nil {
		return err
	}
	fmt.Printf("%v", dump)

	return nil
}
