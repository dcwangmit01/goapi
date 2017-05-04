package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dcwangmit01/goapi/client"
	"github.com/dcwangmit01/goapi/jwt"
	"github.com/dcwangmit01/goapi/util"
)

var (
	optionTokenOnly bool
)

func init() {
	RootCmd.AddCommand(authRootCmd)
	authRootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringVarP(&optionUsername, "username", "u", "", "Username for authentication")
	loginCmd.PersistentFlags().StringVarP(&optionPassword, "password", "p", "", "Password for authentication")
	loginCmd.Flags().BoolVarP(&optionTokenOnly, "token-only", "t", false, "Output only the token")
}

var authRootCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the API service",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Auth and save the JWT token to config",
	Long: `Authenticate against the API auth endpoint.
  * Uses --username and --password if specified
    * Else obtains user input from stdin
  * Hit the /auth endpoint
  * Save the token in the config
  * Print the token to the screen`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return appAuthLogin(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

func appAuthLogin(cmd *cobra.Command, args []string) error {

	// validate args
	if len(args) != 0 {
		cmd.Usage()
		return invalidInputErr
	}

	// get the auth token
	tokenStr, err := client.GetAuthTokenFromOptionOrConfigOrStdin(
		true,           // optionTry
		optionUsername, // optionUsername
		optionPassword, // optionPassword
		false,          // optionContinue
		false,          // configTry
		false,          // configContinue
		true,           // stdinTry
		true,           // stdinContinue
		true,           // saveNewToken
	)
	if err != nil {
		return err
	}

	// print output
	if optionTokenOnly == true {
		// print the token
		fmt.Println(tokenStr)
	} else {
		token, _, err := jwt.ParseJwt(tokenStr)
		if err != nil {
			return err
		}

		dump, err := util.StructToYamlStr(token)
		if err != nil {
			return err
		}
		fmt.Printf("%v", dump)
	}

	return nil
}
