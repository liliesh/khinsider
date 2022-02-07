package main

import (
	"github.com/marcus-crane/khinsider/v2/cmd/khinsider"
)

var (
	version = "v2.0.4"
	commit  = "n/a"
	date    = "n/a"
	builtBy = "go install"
)

func main() {
	buildInfo := khinsider.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}
	khinsider.Execute(buildInfo)
}
