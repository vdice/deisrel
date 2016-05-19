package docker

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"

	git "github.com/deis/deisrel/git"
	reg "github.com/deis/deisrel/registry"
)

const (
	ciRepoOrg    = "deisci"
	prodRepoOrg  = "deis"
	monitorImage = "monitor"
)

// CheckTags checks that all shas provided in repoAndShas exist in tag form
// in provided registries
func CheckTags(ghClient *github.Client, quay reg.Registry, hub reg.Registry, repoAndShas []git.RepoAndSha, ref, newTag string) ([]string, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	resultsCh := make(chan string)
	errCh := make(chan error)
	defer close(errCh)
	for _, rs := range repoAndShas {
		wg.Add(1)
		go func(rs git.RepoAndSha) {
			defer wg.Done()
			var imgTags []reg.ImageAndTag

			// a given GitHub repo may contain multiple components,
			// each publishing images to their own quay repos
			for _, componentName := range repoToComponentNames[rs.RepoName] {
				tag, err := git.ImageTagTransform(rs.Sha)
				if err != nil {
					errCh <- err
				}
				imgTags = append(imgTags,
					reg.ImageAndTag{
						Image: fmt.Sprintf("%s/%s", ciRepoOrg, componentToImageName[componentName]),
						Tag:   tag,
					})
			}
			for _, imgTag := range imgTags {
				if err := quay.CheckExistence(imgTag); err != nil {
					if err != nil {
						errCh <- err
						return
					}
				}
				if err := hub.CheckExistence(imgTag); err != nil {
					if err != nil {
						errCh <- err
						return
					}
				}
				resultsCh <- fmt.Sprintf("%s:%s", imgTag.Image, imgTag.Tag)
			}
		}(rs)
	}
	go func() {
		// wait for all fetches from github to be complete before returning
		wg.Wait()
		close(done)
	}()

	results := []string{}
	errs := []error{}
	for {
		select {
		case <-done:
			return results, errs
		case result := <-resultsCh:
			results = append(results, result)
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
}
