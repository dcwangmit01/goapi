package cmd

import (
	"errors"
	"fmt"

	"github.com/dcwangmit01/goapi/config"
	"github.com/dcwangmit01/goapi/util"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(configRootCmd)
	configRootCmd.AddCommand(configGetCmd)
	configRootCmd.AddCommand(configListCmd)
	configRootCmd.AddCommand(configSetCmd)
	configRootCmd.AddCommand(configUnsetCmd)
}

var configRootCmd = &cobra.Command{
	Use:   "config",
	Short: "Client used to set Key/Value on gRPC service",
	Example: fmt.Sprintf(`  # List all configs:
  %[1]v config list

  # Get a config value:
  %[1]v config get <key>

  # Set a config value:
  %[1]v config set <key> <value>

  # Unset a config value:
  %[1]v config unset <key>`, config.GetAppName()),
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return configList(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	RunE: func(cmd *cobra.Command, args []string) error {
		return configGet(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set new Key/Value on gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return configSet(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Unset new Key/Value on gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return configUnset(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

func configList(cmd *cobra.Command, args []string) error {
	// validate args
	if len(args) != 0 {
		cmd.Usage()
		return invalidInputErr
	}

	// list all of the settings
	yaml, err := util.StructToYamlStr(config.Viper.AllSettings())
	if err != nil {
		return err
	}
	fmt.Printf(yaml)

	return nil
}

func configGet(cmd *cobra.Command, args []string) error {
	// validate args
	if len(args) != 1 {
		cmd.Usage()
		return invalidInputErr
	}

	// print the specific key
	key := args[0]
	if !config.Viper.IsSet(key) {
		return errors.New("config key not found")
	}
	fmt.Printf("%v\n", config.Viper.Get(key))

	return nil
}

func configSet(cmd *cobra.Command, args []string) error {
	// validate args
	if len(args) != 2 {
		cmd.Usage()
		return invalidInputErr
	}

	// save the key/value pair
	key := args[0]
	val := args[1]
	if !config.Viper.IsSet(key) {
		return errors.New("config key not found")
	}
	config.Viper.Set(key, val)
	fmt.Printf("%s=%v\n", key, config.Viper.Get(key))

	// save the config file
	return config.SaveConfig()
}

func configUnset(cmd *cobra.Command, args []string) error {
	// validate args
	if len(args) != 1 {
		cmd.Usage()
		return invalidInputErr
	}

	// unset the value for the key
	key := args[0]
	if !config.Viper.IsSet(key) {
		return errors.New("config key not found")
	}
	config.Viper.Set(key, nil)
	fmt.Printf("%s=%v\n", key, config.Viper.Get(key))

	// save the config file
	return config.SaveConfig()
}
