package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KarpelesLab/rest"
)

var (
	GIT_TAG  string
	DATE_TAG string
)

func main() {
	rest.Host = "www.shells.com"

	// let's make sure we're logged in
	auth, err := checkLogin()
	if err != nil {
		log.Printf("login failed: %s", err)
		os.Exit(1)
	}

	err = run(auth, os.Args[1:])
	if err != nil {
		log.Printf("failed: %s", err)
		os.Exit(1)
	}
}

func showVersion(ri *runInfo) error {
	fmt.Fprintf(os.Stdout, "shells-cli version %s built %s\r\n", GIT_TAG, DATE_TAG)
	return nil
}
