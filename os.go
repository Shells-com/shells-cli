package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func osImgUpload(ri *runInfo) error {
	osId := ri.flags["os"]
	fn := ri.flags["file"]

	fp, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("while trying to open file to upload: %w", err)
	}

	st, err := fp.Stat()
	if err != nil {
		return fmt.Errorf("while grabbing file info: %w", err)
	}

	// we do not implement upload method for that, fail if file is too big
	if st.Size() > 5*1024*1024*1024 {
		return fmt.Errorf("cannot upload files over 5GB")
	}

	// prepare upload
	var upload *fileUpload
	err = ri.auth.Apply(context.Background(), "Shell/OS/"+osId+"/Image:upload", "POST", map[string]interface{}{"filename": filepath.Base(fn)}, &upload)
	if err != nil {
		return err
	}

	log.Printf("Uploading %s (%d bytes) ...", fn, st.Size())

	req, err := http.NewRequest(http.MethodPut, upload.PUT, fp)
	if err != nil {
		return err
	}
	req.ContentLength = st.Size()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to upload: status %s", resp.Status)
	}

	// call completion
	log.Printf("Upload completed, confirming...")
	var img *ShellOsImage
	err = ri.auth.Apply(context.Background(), upload.Complete, "POST", map[string]interface{}{}, &img)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s %s %s\r\n", img.Id, img.Version, img.Filename)

	return nil
}
