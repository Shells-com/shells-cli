package main

import (
	"context"
	"log"
)

type Shell struct {
	Id    string `json:"Shell__"`
	Label string `json:"Label"`
}

func lsShells(auth *authInfo) error {
	// list available shells
	var list []Shell

	err := auth.Apply(context.Background(), "Shell", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, shell := range list {
		log.Printf("%s", shell.Label)
	}

	return nil
}
