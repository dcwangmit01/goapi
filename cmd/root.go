package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	invalidInputErr = errors.New("Invalid Command Line Input")
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "goapi",
	Short: "All-in-one CLI tool for running the goapi server and client",
	Long: `All-in-one CLI tool for running the goapi server and client
  * Source code and documentation is available at
    https://github.com/dcwangmit01/goapi`,
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
