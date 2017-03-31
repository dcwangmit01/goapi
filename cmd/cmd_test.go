package cmd_test

import (
	"github.com/dcwangmit01/grpc-gw-poc/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io"
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

		// set the output to both Stdout and a byteBuffer
		var buf bytes.Buffer
		mw := io.MultiWriter(&buf, os.Stdout)
		cmd.RootCmd.SetOutput(mw)

		// Run the command which parses os.Args
		err := cmd.RootCmd.Execute()

		// restore the output to Stdout
		cmd.RootCmd.SetOutput(os.Stdout)

		// process the output
		//  (?s): allows for "." to represent "\n"
		Expect(buf.String()).Should(MatchRegexp("(?s)grpc-gw-poc.*help.*keyval"))
		Expect(err).Should(BeNil())
	})
})
