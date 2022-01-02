package khinsider

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func (b BuildInfo) Print() {
	fmt.Println("Version:\t", b.Version)
	fmt.Println("Commit Hash:\t", b.Commit)
	fmt.Println("Build Date:\t", b.Date)
	fmt.Println("Build Source:\t", b.BuiltBy)
}

func Execute(buildInfo BuildInfo) {
	cli.VersionPrinter = func(c *cli.Context) {
		buildInfo.Print()
	}
	app := &cli.App{
		Name:     "khinsider",
		Version:  "2.0.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Marcus Crane",
				Email: "khinsider@utf9k.net",
			},
		},
		Usage: "easily fetch videogame soundtracks from downloads.khinsider.com",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				pterm.EnableDebugMessages()
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "search for an album to download",
				Before: func(c *cli.Context) error {
					return BeforeSearch()
				},
				Action: func(c *cli.Context) error {
					return SearchAction()
				},
			},
			{
				Name:    "album",
				Aliases: []string{"a"},
				Usage:   "download an album given a slug",
				Action: func(c *cli.Context) error {
					return DownloadAction(c.Args().First())
				},
			},
			{
				Name:    "index",
				Aliases: []string{"i"},
				Usage:   "generates a local index of all khinsider content",
				Hidden:  true,
				Action: func(c *cli.Context) error {
					return IndexAction()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
