package actions

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmGenerateWorkflow is the cli handler for generating a helm parameters file for deis-workflow
func HelmGenerateWorkflow(ghClient *github.Client) func(*cli.Context) {
	return func(c *cli.Context) {
		log.Printf("Not yet implemented")
	}
}
