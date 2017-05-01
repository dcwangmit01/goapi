package cmd

import (
	"github.com/spf13/cobra"

	svr "github.com/dcwangmit01/goapi/server"

	// do not import, but trigger the init functions which
	// register the grpc service into the serviceregistry
	_ "github.com/dcwangmit01/goapi/example/service"
	_ "github.com/dcwangmit01/goapi/service"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "server",
	Short: "Launches the example webserver on https://localhost:10080",
	Run: func(cmd *cobra.Command, args []string) {
		svr.StartServer()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
