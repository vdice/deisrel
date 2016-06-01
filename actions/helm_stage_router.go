package actions

import (
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmStageRouter is the cli handler for generating a release helm chart for router
func HelmStageRouter(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		helmStage(ghClient, c, RouterChart)
		return nil
	}
}
