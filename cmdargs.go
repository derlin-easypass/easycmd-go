package main

import (
	"flag"
	"fmt"
	"github.com/derlin-easypass/easycmd-go/account"
	"os"
)

var creds *account.Creds = &account.Creds{}
var isNewSession bool
var jsonPath string
var commands string

var remainingArgs []string

func ReadCmd() {
	flag.StringVar(&creds.Password, "pass", "", "password (not recommended !)")
	flag.StringVar(&commands, "cmd", "", "commands (non-interactive execution)")
	flag.StringVar(&jsonPath, "json", "", "Combined with '-new', load session from a json file, which will be encrypted to 'path'.")
	flag.BoolVar(&isNewSession, "new", false, "create a new session saved to 'path'")

	flag.Parse()

	remainingArgs = flag.Args()
	if len(remainingArgs) == 0 {
		fmt.Println("Missing required parameter <path to file>")
		os.Exit(1)
	}

	creds.Path = remainingArgs[0]

	if _, err := os.Stat(creds.Path); os.IsNotExist(err) {
		// session does not exist; ok only if new session flag
		if !isNewSession {
			fmt.Printf("path '%s' does not exist. To create a new session, use '-new' option.\n", creds.Path)
			os.Exit(1)
		}
	} else {
		// session exists. Ensure not to override it if the new flag is used.
		if isNewSession {
			fmt.Printf("path '%s' already exists and the '-new' option is used. Please delete it first or use another path.\n", creds.Path)
			os.Exit(1)
		}

		if jsonPath != "" {
			fmt.Println("warning: -json option without -new. Will be ignored.")
		}
	}

}
