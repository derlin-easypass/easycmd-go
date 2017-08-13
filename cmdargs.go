package main

import (
	"flag"
	"os"
	"fmt"
)

var sessionPath string
var password string
var isNewSession bool
var jsonPath string

var remainingArgs []string

func init() {
	flag.StringVar(&sessionPath, "path", "", "path to the session file")
	flag.StringVar(&password, "pass", "", "password (not recommended !)")	
	flag.StringVar(&jsonPath, "json", "", "Combined with '-new', load session from a json file, which will be encrypted to 'path'.")	
	flag.BoolVar(&isNewSession, "new", false, "create a new session saved to 'path'")

	flag.Parse()

	remainingArgs = flag.Args()
	fmt.Println(remainingArgs)

	if sessionPath == "" {
		fmt.Println("Missing required parameter -path <path to file>")
		os.Exit(1)
	}

	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		// session does not exist; ok only if new session flag
		if !isNewSession {
			fmt.Printf("path '%s' does not exist. To create a new session, use '-new' option.\n", sessionPath)
			os.Exit(1)
		}
	}else{
		// session exists. Ensure not to override it if the new flag is used.
		if isNewSession {
			fmt.Printf("path '%s' already exists and the '-new' option is used. Please delete it first or use another path.\n", sessionPath)
			os.Exit(1)
		}

		if jsonPath != "" {
			fmt.Println("warning: -json option without -new. Will be ignored.")
		}
	}

}