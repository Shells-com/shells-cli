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

type ShellOsImage struct {
	Id       string `json:"Shell_OS_Image__"`
	Version  string
	QA       string `json:"QA_Passed"` // Y, P or N
	Filename string
	Format   string
	Source   string
	Status   string
	Size     string // as string because might be too large to fit
	Hash     string
	// Created timestamp
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

func osImgList(ri *runInfo) error {
	// list available shells
	var list []ShellOsImage

	osId := ri.flags["os"]

	err := ri.auth.Apply(context.Background(), "Shell/OS/"+osId+"/Image", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, img := range list {
		fmt.Fprintf(os.Stdout, "%s %s %s\r\n", img.Id, img.Version, img.Filename)
	}

	return nil
}
