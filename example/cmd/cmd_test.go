package cmd_test

import (
	"github.com/dcwangmit01/goapi/config"
	"github.com/dcwangmit01/goapi/example/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io"
	"os"
)

var debug = false

var _ = Describe("Cmd", func() {
	// The Ginkgo test runner takes over os.Args and fills it with its own
	// flags.  This makes the cobra command arg parsing fail because of
	// unexpected options.  Work around this.

	// Save a copy of os.Args
	var origArgs = os.Args[:]

	// A non-threadsafe buffer for capturing stdout
	var buf bytes.Buffer

	BeforeEach(func() {
		// Trim os.Args to only the first arg, which is the command itself
		os.Args = os.Args[:1]

		// set the output to both Stdout and a byteBuffer
		var mw io.Writer
		if debug == true {
			mw = io.MultiWriter(&buf, os.Stdout)
		} else {
			mw = io.MultiWriter(&buf)
		}
		cmd.RootCmd.SetOutput(mw)
	})

	AfterEach(func() {
		// Restore os.Args
		os.Args = origArgs[:]

		// restore the output to Stdout
		cmd.RootCmd.SetOutput(os.Stdout)
	})

	Describe("RootCmd", func() {
		Context("When run with no args", func() {
			It("Should show rootcmd help", func() {
				// Run the command which outputs to stdout
				err := cmd.RootCmd.Execute()

				// bytes.Buffer.String() returns the contents of the unread
				// portion of the buffer as a string.
				out := buf.String()

				// process the output
				//  (?s): allows for "." to represent "\n"
				Expect(out).Should(MatchRegexp(format.Sprintf("(?s)%v.*help.*keyval", config.GetAppName())))
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("keyval Subcommand", func() {
		Context("When run with no args", func() {
			It("Should show keyval help", func() {
				// Set args to command
				os.Args = append(os.Args, "keyval")

				// Run the command which outputs to stdout
				err := cmd.RootCmd.Execute()

				// bytes.Buffer.String() returns the contents of the unread
				// portion of the buffer as a string.
				out := buf.String()

				// process the output
				//  (?s): allows for "." to represent "\n"
				Expect(out).Should(MatchRegexp(format.Sprintf("(?s)%v.*keyval.*create.*read.*update.*delete", config.GetAppName())))
				Expect(err).Should(BeNil())
			})
		})
	})
})
