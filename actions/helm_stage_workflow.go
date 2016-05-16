package actions

import (
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmStageWorkflow is the cli handler for generating a release helm chart for workflow
func HelmStageWorkflow(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		helmStage(ghClient, c, WorkflowChart)
		return nil
	}
}
