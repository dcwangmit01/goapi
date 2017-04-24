package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	log "github.com/Sirupsen/logrus"
	"github.com/dcwangmit01/goapi/app/logutil"

	pb "github.com/dcwangmit01/goapi/app"
	config "github.com/dcwangmit01/goapi/app/config"
	certs "github.com/dcwangmit01/goapi/resources/certs"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "goapi",
	Short: "All-in-one command for running GRPC server, gw, and client",
	Long: `
Run the combined GRPC server and gateway

    goapi serve

Run the GRPC client

    goapi client <strings to be echo'd ...>
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goapi.yaml)")

	// Cobra also supports local flags, which will only run when
	// this action is called directly.  Any flags ending with "P"
	// adds short single character option.
	// RootCmd.Flags().BoolP("foo", "f", false, "Help message for foo")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".goapi") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Helper method used by many commands to dial to the GRPC server, and
// then run a callback function immediately after connection.
func grpcDialAndRunAuth(callback_func func(pb.AuthClient, context.Context)) {

	// todo: send a token instead of a test header
	md := metadata.Pairs("timestamp", time.Now().Format(time.RFC3339))
	ctx := metadata.NewContext(context.Background(), md)

	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certs.CertPool, config.ServerAddress)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(config.ServerAddress, opts...)
	if err != nil {
		logutil.AddCtx(log.WithFields(log.Fields{
			"error": err,
		})).Info("Failed to Dial")
	}
	defer conn.Close()

	//grpc.SendHeader(ctx, header)
	client := pb.NewAuthClient(conn)

	callback_func(client, ctx)
}

// Helper method used by many commands to dial to the GRPC server, and
// then run a callback function immediately after connection.
func grpcDialAndRunKeyVal(callback_func func(pb.KeyValClient, context.Context)) {

	// todo: send a token instead of a test header
	md := metadata.Pairs("timestamp", time.Now().Format(time.RFC3339))
	ctx := metadata.NewContext(context.Background(), md)

	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certs.CertPool, config.ServerAddress)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(config.ServerAddress, opts...)
	if err != nil {
		logutil.AddCtx(log.WithFields(log.Fields{
			"error": err,
		})).Info("Failed to Dial")
	}
	defer conn.Close()

	//grpc.SendHeader(ctx, header)
	client := pb.NewKeyValClient(conn)

	callback_func(client, ctx)
}
