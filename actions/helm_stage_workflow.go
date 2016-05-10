package actions

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmStageWorkflow is the cli handler for generating a release helm chart for workflow
func HelmStageWorkflow(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		const chartDir = "workflow-dev"
		var fileNames = []string{
			fmt.Sprintf("%s/README.md", chartDir),
			fmt.Sprintf("%s/Chart.yaml", chartDir),
		}

		stagingDir := filepath.Join(stagingPath, chartDir)
		if err := createDir(ourFS, stagingDir); err != nil {
			log.Fatalf("Error creating dir %s (%s)", stagingDir, err)
		}
		helmStage(ghClient, c, fileNames, stagingDir)

		return nil
	}
}
