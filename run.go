package main

import (
	"log"
	"strings"
)

type runInfo struct {
	auth  *authInfo
	args  []string
	flags map[string]string
}

func run(auth *authInfo, args []string) error {
	ri := &runInfo{
		auth:  auth,
		args:  args,
		flags: make(map[string]string),
	}

	return ri.handle(rootCmd)
}

func (ri *runInfo) handle(cmd *cmdInfo) error {
	for {
		if cmd.cb != nil {
			return cmd.cb(ri)
		}
		if len(ri.args) > 0 {
			if v, ok := cmd.children[strings.ToLower(ri.args[0])]; ok {
				ri.args = ri.args[1:]
				cmd = v
				continue
			}
		} else {
			break
		}
	}

	log.Printf("Error: invalid argument provided")
	log.Printf("Please choose one of:")

	for k := range cmd.children {
		log.Printf(" * %s", k)
	}
	return nil
}
