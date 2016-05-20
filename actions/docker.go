package actions

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/coreos/go-quay/quay"
	"github.com/google/go-github/github"

	reg "github.com/deis/deisrel/registry"
)

const (
	ciRepoOrg    = "deisci"
	prodRepoOrg  = "deis"
	monitorImage = "monitor"
)

// TODO: place in 'docker' package, https://github.com/deis/deisrel/issues/60
func dockerCheckTags(ghClient *github.Client, quay reg.Registry, hub reg.Registry, repoAndShas []repoAndSha, ref string) ([]string, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	resultsCh := make(chan string)
	errCh := make(chan error)
	defer close(errCh)
	for _, rs := range repoAndShas {
		wg.Add(1)
		go func(rs repoAndSha) {
			defer wg.Done()
			var imgTags []reg.ImageAndTag

			// a given GitHub repo may contain multiple components,
			// each publishing images to their own quay repos
			for _, componentName := range repoToComponentNames[rs.repoName] {
				imgTags = append(imgTags,
					reg.ImageAndTag{
						Image: fmt.Sprintf("%s/%s", ciRepoOrg, componentToImageName[componentName]),
						Tag:   imageTagTransform(rs.sha),
					})
			}
			for _, imgTag := range imgTags {
				if err := quay.CheckExistence(imgTag); err != nil {
					errCh <- err
					return
				}
				if err := hub.CheckExistence(imgTag); err != nil {
					errCh <- err
					return
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

// TODO: refactor to just add `-push` flag to check-tags
// maybe rename check-tags back to retag and supply the newTag no matter what
func dockerPushTags(ghClient *github.Client, quay reg.Registry, hub reg.Registry, repoAndShas []repoAndSha, ref, newTag string) ([]string, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	resultsCh := make(chan string)
	errCh := make(chan error)
	defer close(errCh)
	// TODO: remove!
	repoAndShas = []repoAndSha{
		repoAndSha{repoName: "controller", sha: "9979523faf255c6123c430d1afcda858602993c6"},
	}
	for _, rs := range repoAndShas {
		wg.Add(1)
		go func(rs repoAndSha) {
			defer wg.Done()
			imgRetags := make(map[reg.ImageAndTag]reg.ImageAndTag)

			// a given GitHub repo may contain multiple components,
			// each publishing images to their own quay repos
			for _, componentName := range repoToComponentNames[rs.repoName] {
				imgName := componentToImageName[componentName]
				origImgTag := reg.ImageAndTag{
					Image: fmt.Sprintf("%s/%s", ciRepoOrg, imgName),
					Tag:   imageTagTransform(rs.sha),
				}
				newImgTag := reg.ImageAndTag{
					// Image: fmt.Sprintf("%s/%s", prodRepoOrg, imgName),
					// TODO: change back to above
					Image: fmt.Sprintf("%s/%s", "vdice", imgName),
					Tag:   newTag,
				}
				imgRetags[origImgTag] = newImgTag
			}
			for orig, new := range imgRetags {
				fmt.Println("TODO: currently not pushing to quay!")
				// if err := quay.PushTag(orig, new); err != nil {
				// 		errCh <- err
				// 		return
				// }
				if err := hub.PushTag(orig, new); err != nil {
					errCh <- err
					return
				}
				resultsCh <- fmt.Sprintf("%s:%s retagged as %s:%s", orig.Image, orig.Tag, new.Image, new.Tag)
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

func setupRegistries() (reg.Registry, reg.Registry, error) {
	hubUsername := os.Getenv("DOCKER_USER")
	hubPassword := os.Getenv("DOCKER_PASSWORD")
	if hubUsername == "" || hubPassword == "" {
		return nil, nil, errors.New("DOCKER_USER and DOCKER_PASSWORD must be set to use DockerHub registry")
	}
	hubClient, err := reg.GetHub("https://index.docker.io/", hubUsername, hubPassword)
	if err != nil {
		return nil, nil, err
	}
	hub := reg.NewHubRegistry(hubClient)

	quay := reg.NewQuayRegistry(quay.NewHTTPClient(nil), reg.NewQuayAuth())

	return hub, quay, nil
}
