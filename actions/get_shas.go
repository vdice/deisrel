package actions

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

const (
	// ShortFlag is the cli flag that indicates whether to show short or long SHAs
	ShortFlag = "short"
)

// GetShas is the CLI action for getting github shas of all of the Deis Workflow repos
func GetShas(ghClient *github.Client) func(c *cli.Context) {
	return func(c *cli.Context) {
		transformFunc := func(s string) string { return s }
		if c.Bool(ShortFlag) {
			transformFunc = func(s string) string { return s[:7] }
		}
		shas, err := getShas(ghClient, repoNames, transformFunc)
		if err != nil {
			log.Fatal(err)
		}
		for _, sha := range shas {
			fmt.Print(sha)
		}
	}
}
