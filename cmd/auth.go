package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dcwangmit01/goapi/app/client"
	"github.com/dcwangmit01/goapi/app/jwt"
)

var (
	tokenOnly bool
)

func init() {
	RootCmd.AddCommand(authRootCmd)
	authRootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolVarP(&tokenOnly, "token-only", "t", false, "Output only the token")

}

var authRootCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with goapi service",
	Long:  ``,
}

var loginCmd = &cobra.Command{
	Use:   "login USERNAME PASSWORD",
	Short: "Auth and save the JWT token to config",
	Long: `Authenticate against the API auth endpoint.
  * With the provided USERNAME AND PASSWORD
  * Hit the /auth endpoint
  * Save the token in the config
  * Print the token to the screen`,
	Example: `asdf`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return appAuthLogin(cmd, args)
	},
	SilenceUsage: true, // mark true otherwise usage is printed on EVERY error
}

func appAuthLogin(cmd *cobra.Command, args []string) error {

	// validate args
	if len(args) != 2 {
		cmd.Usage()
		return invalidInputErr
	}

	// authenticate
	tokenStr, err := client.Authenticate(
		args[0], // username
		args[1], // password
	)
	if err != nil {
		return err
	}

	if tokenOnly == true {
		// print the token
		fmt.Println(tokenStr)
	} else {
		token, _, err := jwt.ParseJwt(tokenStr)
		if err != nil {
			return err
		}

		dump, err := client.StructToYamlStr(token)
		if err != nil {
			return err
		}

		fmt.Printf("%v\n", dump)
	}
	return nil
}
