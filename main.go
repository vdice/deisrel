package main

import (
	"log"
	"os"

	"github.com/arschles/deisrel/actions"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ghTkn := os.Getenv("GITHUB_TOKEN")
	if ghTkn == "" {
		log.Fatalf("'GITHUB_TOKEN' env var required")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghTkn})
	cl := oauth2.NewClient(oauth2.NoContext, ts)
	ghClient := github.NewClient(cl)

	app := cli.NewApp()
	app.Name = "deisrel"
	app.Usage = "Utilities for releasing a new Deis version"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "shas",
			Action: actions.GetShas(ghClient),
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
