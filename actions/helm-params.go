package actions

import (
	// "fmt"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

const (
	// TagFlag represents the '-tag' flag
	TagFlag = "tag"
	// PullPolicyFlag represents the '-pull-policy' flag
	PullPolicyFlag = "pull-policy"
	// OrgFlag represents the '-org' flag
	OrgFlag = "org"
)

// HelmParams is the cli handler for generating a helm parameters file
func HelmParams(ghClient *github.Client) func(c *cli.Context) {
	return func(c *cli.Context) {
		// tag := c.StringFlag(TagFlag)
		// pullPolicy := c.StringFlag(PullPolicyFlag)
		// org := c.StringFlag(OrgFlag)

	}
}
