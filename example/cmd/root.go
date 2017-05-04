package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dcwangmit01/goapi/config"
	"github.com/spf13/cobra"
)

var (
	optionUsername string
	optionPassword string

	optionConfigFile string
	invalidInputErr  = errors.New("Invalid Command Line Input")
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   config.GetAppName(),
	Short: fmt.Sprintf("All-in-one CLI tool for running the %v server and client", config.GetAppName()),
	Long:  fmt.Sprintf("All-in-one CLI tool for running the %v server and client", config.GetAppName()),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// The cobra framework will print the error automatically, so
		// no need to do that here.
		// fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(-1)
	}
}

/*
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	//RootCmd.PersistentFlags().StringVar(&optionConfigFile, "config", "",
	//	fmt.Sprintf("config file (default is $HOME/.%v/config.yaml)", config.GetAppName()))

	// Cobra also supports local flags, which will only run when
	// this action is called directly.  Any flags ending with "P"
	// adds short single character option.
	// RootCmd.Flags().BoolP("foo", "f", false, "Help message for foo")
}
*/
