package actions

import (
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmStageE2E is the cli handler for generating a release helm chart for deis-e2e
func HelmStageE2E(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		helmStage(ghClient, c, WorkflowE2EChart)
		return nil
	}
}
