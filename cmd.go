package main

type cmdInfo struct {
	children map[string]*cmdInfo
	cb       func(*runInfo) error
	flags    []cmdFlag
}

type cmdFlag struct {
	Name     string // name as on command line
	Usage    string
	Required bool // if true, value is required
}

var rootCmd = &cmdInfo{
	children: map[string]*cmdInfo{
		"version": {cb: showVersion},
		"shells": {
			children: map[string]*cmdInfo{
				"ls":    {cb: shellsList},
				"info":  {cb: shellsInfo, flags: []cmdFlag{{Name: "shell", Usage: "Specify shell to fetch", Required: true}}},
				"start": {cb: shellsStart, flags: []cmdFlag{{Name: "shell", Usage: "Specify shell to start", Required: true}}},
				"view":  {cb: shellsView, flags: []cmdFlag{{Name: "shell", Usage: "Specify shell to get view URL", Required: true}}},
			},
		},
		"os": {
			children: map[string]*cmdInfo{
				"image": {
					children: map[string]*cmdInfo{
						"ls": {cb: osImgList, flags: []cmdFlag{{Name: "os", Usage: "Specify the OS to list images", Required: true}}},
						"upload": {
							cb: osImgUpload,
							flags: []cmdFlag{
								{Name: "os", Usage: "Specify the OS to upload to", Required: true},
								{Name: "file", Usage: "File to be uploaded", Required: true},
							},
						},
					},
				},
				"ls": {cb: osList},
			},
		},
	},
}
