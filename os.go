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

func osList(ri *runInfo) error {
	// list available shells
	var list []ShellOs

	err := ri.auth.Apply(context.Background(), "Shell/OS", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, shos := range list {
		fmt.Fprintf(os.Stdout, "%s %s\r\n", shos.Id, shos.Name)
	}

	return nil
}
