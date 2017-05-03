package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func CredentialsFromStdin() (string, string, error) {
	// Modified from:
	//   http://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go

	reader := bufio.NewReader(os.Stdin)

	// read username
	fmt.Fprintf(os.Stderr, "Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	// read password
	fmt.Fprintf(os.Stderr, "Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)

	fmt.Fprintf(os.Stderr, "\n")
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
