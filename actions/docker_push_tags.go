package actions

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// DockerPushTags is the CLI action for checking that Docker image tags
// exist in registries we are interested in
func DockerPushTags(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ref := c.GlobalString(RefFlag)
		newTag := c.Args().Get(0)
		if newTag == "" {
			log.Fatalln("Usage: docker push-tags <new-tag>")
		}

		hub, quay, err := setupRegistries()
		if err != nil {
			log.Fatalf("Error encountered setting up registries (%s)", err)
		}

		// latest sha values should exist in 'git-<sha>' form as tags in registries
		repoAndShas, err := getShas(ghClient, repoNames, shortShaTransform, ref)
		if err != nil {
			return err
		}

		pushedImgTags, errs := dockerPushTags(ghClient, quay, hub, repoAndShas, ref, newTag)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Printf("Error encountered attempting pushing tag (%s)", err)
			}
		}

		fmt.Println("Successfully pushed all image tags to Quay.io and DockerHub registries:")
		for _, pushedImgTag := range pushedImgTags {
			fmt.Println(pushedImgTag)
		}
		return nil
	}
}
