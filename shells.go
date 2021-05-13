package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/KarpelesLab/rest"
)

type Shell struct {
	Id         string `json:"Shell__"`
	Label      string
	Egine      string // "full"
	Size       int
	Status     string
	State      string
	SshPort    int `json:"Ssh_Port"`
	Username   string
	Hostname   string
	MAC        string
	IPv4       string // internal ipv4
	IPv6       string // internal ipv6
	Created    rest.Time
	Expires    rest.Time
	Host       *ShellHost
	IPs        []*ShellIP
	OS         *ShellOs
	Datacenter *ShellDatacenter `json:"Shell_Datacenter"`

	// Ephemeral_Viewer
	EphemeralViewer string `json:"Ephemeral_Viewer"`
}

type ShellHost struct {
	Id       string `json:"Shell_Host__"`
	Name     string // "ams01-09"
	IP       string
	PublicIP string `json:"Public_IP"`
	IPv6     string
	Kernel   string // host kernel
}

type ShellIP struct {
	IP     string
	Family string // "ipv4"
	Type   string // "nat", "route" or "anycast"
}

type ShellDatacenter struct {
	Name     string // "ams01"
	Location string // "Amsterdam (Europe)"
	Country  string `json:"Country__"`
}

func shellsList(ri *runInfo) error {
	// list available shells
	var list []Shell

	err := ri.auth.Apply(context.Background(), "Shell", "GET", map[string]interface{}{}, &list)
	if err != nil {
		return err
	}

	for _, shell := range list {
		fmt.Fprintf(os.Stdout, "%s %s (%s)\r\n", shell.Id, shell.Label, shell.State)
	}

	return nil
}

func shellsInfo(ri *runInfo) error {
	var shell *Shell

	err := ri.auth.Apply(context.Background(), "Shell/"+ri.flags["shell"], "GET", map[string]interface{}{}, &shell)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Id: %s\r\n", shell.Id)
	fmt.Fprintf(os.Stdout, "Label: %s\r\n", shell.Label)
	fmt.Fprintf(os.Stdout, "Size: %d units\r\n", shell.Size)
	fmt.Fprintf(os.Stdout, "State: %s\r\n", shell.State)
	fmt.Fprintf(os.Stdout, "Host: %s\r\n", shell.Host.Name)

	return nil
}

func shellsState(ri *runInfo) error {
	var shell *Shell

	err := ri.auth.Apply(context.Background(), "Shell/"+ri.flags["shell"], "GET", map[string]interface{}{}, &shell)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\r\n", shell.State)

	return nil
}

func shellsView(ri *runInfo) error {
	var res linkAccess

	err := ri.auth.Apply(context.Background(), "MetaObject/"+ri.flags["shell"]+":linkAccess", "POST", map[string]interface{}{"access": "A", "expires": "+2 hour"}, &res)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "https://view.shells.com/?id=%s\r\n", url.QueryEscape(res.Link))

	return nil
}

func shellsStart(ri *runInfo) error {
	var res map[string]interface{}

	err := ri.auth.Apply(context.Background(), "Shell/"+ri.flags["shell"]+":start", "POST", map[string]interface{}{}, &res)
	if err != nil {
		return err
	}

	return nil
}

func shellsStop(ri *runInfo) error {
	var res map[string]interface{}

	err := ri.auth.Apply(context.Background(), "Shell/"+ri.flags["shell"]+":stop", "POST", map[string]interface{}{}, &res)
	if err != nil {
		return err
	}

	return nil
}

func shellsRestart(ri *runInfo) error {
	var res map[string]interface{}

	err := ri.auth.Apply(context.Background(), "Shell/"+ri.flags["shell"]+":restart", "POST", map[string]interface{}{}, &res)
	if err != nil {
		return err
	}

	return nil
}

func shellsReboot(ri *runInfo) error {
	var res map[string]interface{}

	err := ri.auth.Apply(context.Background(), "Shell/"+ri.flags["shell"]+":reboot", "POST", map[string]interface{}{}, &res)
	if err != nil {
		return err
	}

	return nil
}
