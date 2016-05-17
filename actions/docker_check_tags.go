package actions

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-quay/quay"
	"github.com/google/go-github/github"

	reg "github.com/deis/deisrel/registry"
)

// DockerCheckTags is the CLI action for checking that Docker image tags
// exist in registries we are interested in
func DockerCheckTags(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ref := c.GlobalString(RefFlag)

		// setup registries
		hubClient, err := reg.GetHub("https://index.docker.io/", "", "")
		if err != nil {
			return err
		}
		hub := reg.NewHubRegistry(hubClient)
		quay := reg.NewQuayRegistry(quay.NewHTTPClient(nil), reg.NewQuayAuth())

		// latest sha values should exist in 'git-<sha>' form as tags in registries
		repoAndShas, err := getShas(ghClient, repoNames, shortShaTransform, ref)
		if err != nil {
			return err
		}

		foundImgTags, errs := dockerCheckTags(ghClient, quay, hub, repoAndShas, ref)
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error encountered attempting checking tag (%s)", err)
			}
		}

		log.Println("Successfully found all image tags in Quay.io and DockerHub registries:")
		for _, foundImgTag := range foundImgTags {
			log.Println(foundImgTag)
		}
		return nil
	}
}
