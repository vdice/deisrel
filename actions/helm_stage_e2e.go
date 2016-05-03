package actions

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmStageE2E is the cli handler for generating a release helm chart for deis-e2e
func HelmStageE2E(ghClient *github.Client) func(*cli.Context) {
	return func(c *cli.Context) {
		const chartDir = "workflow-dev-e2e"
		var fileNames = []string{
			fmt.Sprintf("%s/README.md", chartDir),
			fmt.Sprintf("%s/Chart.yaml", chartDir),
		}

		stagingDir := filepath.Join(stagingPath, chartDir)
		if err := createDir(ourFS, stagingDir); err != nil {
			log.Fatalf("Error creating dir %s (%s)", stagingDir, err)
		}
		helmStage(ghClient, c, fileNames, stagingDir)
	}
}
