package main

type cmdInfo struct {
	children map[string]*cmdInfo
	cb       func(*runInfo) error
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
