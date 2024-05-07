package env

import "github.com/urfave/cli/v2"

var appFlags *AppFlags

type AppFlags struct {
	FlacMode   bool
	LocalIndex bool
}

func GetAppFlags() *AppFlags {
	if appFlags == nil {
		appFlags = &AppFlags{}
	}

	return appFlags
}

func SetAppFlags(c *cli.Context) {
	appFlags.FlacMode = c.Bool("flac-mode")
	appFlags.LocalIndex = c.Bool("local-index")
}
