package main

import (
	"context"
	"fmt"
	"os"
)

type Shell struct {
	Id    string `json:"Shell__"`
	Label string `json:"Label"`
}

func shellsList(auth *authInfo) error {
	// list available shells
	var list []Shell

	err := auth.Apply(context.Background(), "Shell", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, shell := range list {
		fmt.Fprintf(os.Stdout, "%s\r\n", shell.Label)
	}

	return nil
}
