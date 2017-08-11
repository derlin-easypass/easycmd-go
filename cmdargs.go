package main

import (
	"flag"
	"os"
	"fmt"
)

var sessionPath string
var password string

func init() {
	flag.StringVar(&sessionPath, "path", "", "path to the session file")
	flag.StringVar(&password, "pass", "", "password (not recommended !)")	

	flag.Parse()

	if sessionPath == "" {
		fmt.Println("Missing required parameter -path <path to file>")
		os.Exit(1)
	}
}