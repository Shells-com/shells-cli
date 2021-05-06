package main

import (
	"log"
	"strings"
)

type cmdInfo struct {
	children map[string]*cmdInfo
	cb       func(*authInfo) error
}

var rootCmd = &cmdInfo{
	children: map[string]*cmdInfo{
		"shells": {
			children: map[string]*cmdInfo{
				"ls": {cb: shellsList},
			},
		},
		"os": {
			children: map[string]*cmdInfo{
				"ls": {cb: osList},
			},
		},
	},
}

func (i *cmdInfo) handle(auth *authInfo, args []string) error {
	if len(args) == 0 {
		if i.cb != nil {
			return i.cb(auth)
		}
	} else if i.children != nil {
		if v, ok := i.children[strings.ToLower(args[0])]; ok {
			return v.handle(auth, args[1:])
		}
	}

	log.Printf("Error: argument required")
	log.Printf("Please choose one of:")

	for k := range i.children {
		log.Printf(" * %s", k)
	}
	return nil
}
