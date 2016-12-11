package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/joshwget/strato/cmd/add"
	"github.com/joshwget/strato/cmd/inspect"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = os.Args[0]
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "verbose",
		},
		cli.StringFlag{
			Name:  "registry",
			Value: "https://registry-1.docker.io/",
		},
		cli.StringFlag{
			Name:  "user",
			Value: "joshwget",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("verbose") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:     "add",
			HideHelp: true,
			Action:   add.Action,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Value: "/",
				},
				cli.StringFlag{
					Name: "skip",
					// TODO: this is a hack to fix e2fsprogs
					Value: "usr/share/man*",
				},
			},
		},
		{
			Name:            "inspect",
			HideHelp:        true,
			SkipFlagParsing: true,
			Action:          inspect.Action,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
