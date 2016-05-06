package actions

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmGenerateE2E is the cli handler for generating a helm parameters file for deis-e2e
func HelmGenerateE2E(ghClient *github.Client) func(*cli.Context) {
	return func(c *cli.Context) {
		params := genParamsComponentAttrs{
			Org:        c.GlobalString(OrgFlag),
			PullPolicy: c.GlobalString(PullPolicyFlag),
			Tag:        c.GlobalString(TagFlag),
		}
		if params.Tag == "" {
			reposAndShas, err := getShas(ghClient, []string{"workflow-e2e"}, shortShaTransform)
			if err != nil {
				log.Fatalf("No tag given and couldn't fetch sha from GitHub (%s)", err)
			} else if len(reposAndShas) < 1 {
				log.Fatalf("No tag given and no sha returned from GitHub for deis/workflow-e2e")
			}
			params.Tag = "git-" + reposAndShas[0].sha
		}
		if err := generateParamsE2ETpl.Execute(os.Stdout, params); err != nil {
			log.Fatalf("Error outputting the e2e values file (%s)", err)
		}
	}
}
