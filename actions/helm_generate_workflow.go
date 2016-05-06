package actions

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmGenerateWorkflow is the cli handler for generating a helm parameters file for deis-workflow
func HelmGenerateWorkflow(ghClient *github.Client) func(*cli.Context) {
	return func(c *cli.Context) {
		const (
			chartDir = "workflow-dev"
		)
		defaultParamsComponentAttrs := genParamsComponentAttrs{
			Org:        c.GlobalString(OrgFlag),
			PullPolicy: c.GlobalString(PullPolicyFlag),
			Tag:        c.GlobalString(TagFlag),
		}
		paramsComponentMap := createParamsComponentMap()
		// fill in map with default values
		for _, componentName := range componentNames {
			paramsComponentMap[componentName] = defaultParamsComponentAttrs
		}

		if c.GlobalString(TagFlag) == "" {
			// gather latest sha for each repo via GitHub api
			reposAndShas, err := getShas(ghClient, repoNames, shortShaTransform)
			if err != nil {
				log.Fatalf("No tag given and couldn't fetch sha from GitHub (%s)", err)
			} else if len(reposAndShas) < 1 {
				log.Fatalf("No tag given and no shas returned from GitHub for %s", defaultParamsComponentAttrs.Org)
			}

			// a given repo may track multiple components; update each component Tag accordingly
			for _, repoAndSha := range reposAndShas {
				repoComponentNames := repoToComponentNames[repoAndSha.repoName]
				paramsComponentAttrs := defaultParamsComponentAttrs
				for _, componentName := range repoComponentNames {
					paramsComponentAttrs.Tag = "git-" + repoAndSha.sha
					paramsComponentMap[componentName] = paramsComponentAttrs
				}
			}
		}

		shouldStage := c.GlobalBool(StageFlag)
		if err := generateParams(shouldStage, ourFS, chartDir, paramsComponentMap); err != nil {
			log.Fatalf("Error outputting the workflow values file (%s)", err)
		}
	}
}
