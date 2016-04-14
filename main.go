package main

import (
	"os"

	"github.com/arschles/deisrel/actions"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "deisrel"
	app.Usage = "Utilities for releasing a new Deis version"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "shas",
			Action: actions.GetShas,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  actions.ShortFlag,
					Usage: "Whether to show short 7 character shas",
				},
			},
		},
	}

	app.Run(os.Args)
}
