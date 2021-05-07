package main

import (
	"fmt"
	"os"
	"sort"
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
	var prefix []string

	for {
		// check flags
		for _, flag := range cmd.flags {
			flagN := "-" + flag.Name
			found := false
			// look in args for this
			for p, arg := range ri.args {
				if strings.ToLower(arg) == flagN {
					// we have something!
					if len(ri.args) < p+2 {
						return fmt.Errorf("Flag %s needs to be followed by an argument: %s", flagN, flag.Usage)
					}
					ri.flags[flag.Name] = ri.args[p+1]
					ri.args = append(ri.args[:p], ri.args[p+2:]...)
					found = true
					break
				}
			}
			if found {
				continue
			}
			if flag.Required {
				return fmt.Errorf("Flag %s is required: %s", flagN, flag.Usage)
			}
		}

		if cmd.cb != nil {
			return cmd.cb(ri)
		}
		if len(ri.args) > 0 {
			if v, ok := cmd.children[strings.ToLower(ri.args[0])]; ok {
				prefix = append(prefix, ri.args[0])
				ri.args = ri.args[1:]
				cmd = v
				continue
			} else {
				break
			}
		} else {
			break
		}
	}

	if len(cmd.children) == 0 {
		return fmt.Errorf("invalid argument provided")
	}

	fmt.Fprintf(os.Stderr, "Error: invalid argument provided\r\n")
	fmt.Fprintf(os.Stderr, "Please choose one of:\r\n")

	list := make([]string, 0, len(cmd.children))
	for k := range cmd.children {
		list = append(list, k)
	}
	sort.Strings(list)

	for _, k := range list {
		fmt.Fprintf(os.Stderr, " * %s\r\n", strings.Join(append(prefix, k), " "))
	}

	os.Exit(2)
	return nil
}
