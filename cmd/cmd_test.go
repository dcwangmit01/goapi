package cmd_test

import (
	"github.com/dcwangmit01/grpc-gw-poc/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
)

var _ = Describe("RootCmd", func() {
	// The Ginkgo test runner takes over os.Args and fills it with its own
	// flags.  This makes the cobra command arg parsing fail because of
	// unexpected options.  Work around this.

	// Save a copy of os.Args
	var origArgs = os.Args[:]

	BeforeEach(func() {
		// Trim os.Args to only the first arg, which is the command itself
		os.Args = os.Args[:1]
	})

	AfterEach(func() {
		// Restore os.Args
		os.Args = origArgs[:]
	})

	It("Should run without arguments", func() {

		// Run the command which parses os.Args
		err := cmd.RootCmd.Execute()

		Expect(err).Should(BeNil())
	})
})
