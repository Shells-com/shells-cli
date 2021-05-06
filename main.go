package main

import (
	"log"
	"os"

	"github.com/KarpelesLab/goupd"
	"github.com/KarpelesLab/rest"
)

func main() {
	log.Printf("shells-cli version %s", goupd.GIT_TAG)
	rest.Host = "www.shells.com"

	// let's make sure we're logged in
	auth, err := checkLogin()
	if err != nil {
		log.Printf("login failed: %s", err)
		os.Exit(1)
	}

	err = rootCmd.handle(auth, os.Args[1:])
	if err != nil {
		log.Printf("failed: %s", err)
		os.Exit(1)
	}
}
