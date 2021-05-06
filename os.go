package main

import (
	"context"
	"fmt"
	"os"
)

type ShellOs struct {
	Id   string `json:"Shell_OS__"`
	Name string

	URL     string
	Default string   // Y|N
	Ready   string   // Y|N
	Visible string   // Y|N
	Beta    string   // Y|N
	Public  string   // Y|N
	Family  string   // linux|windows|macos|android|unknown
	Boot    string   // guest-linux|bios|efi
	CPU     string   // x86_64
	Purpose string   // unknown|desktop|server|mobile
	Cmdline string   // cmdline for guest-linux
	Flags   []string // byol_warning
}

func osList(auth *authInfo) error {
	// list available shells
	var list []ShellOs

	err := auth.Apply(context.Background(), "Shell/OS", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, shos := range list {
		fmt.Fprintf(os.Stdout, "%s\r\n", shos.Name)
	}

	return nil
}
