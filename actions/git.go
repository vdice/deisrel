package actions

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"
)

func getShas(ghClient *github.Client, repos []string, transform func(string) string) ([]string, error) {
	outCh := make(chan string)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	for _, repo := range repoNames {
		wg.Add(1)
		ch := make(chan string)
		ech := make(chan error)
		go func(repo string) {
			defer wg.Done()
			repoCommits, _, err := ghClient.Repositories.ListCommits("deis", repo, &github.CommitsListOptions{
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			})
			if err != nil {
				ech <- fmt.Errorf("Error listing commits for repo %s (%s)", repo, err)
				return
			}
			if len(repoCommits) < 1 {
				ech <- fmt.Errorf("No commits found for repo %s", repo)
				return
			}
			repoCommit := repoCommits[0]
			sha := transform(*repoCommit.SHA)
			ch <- fmt.Sprintf("%s: %s\n", repo, sha)
		}(repo)
		go func() {
			select {
			case e := <-ech:
				errCh <- e
			case o := <-ch:
				outCh <- o
			}
		}()
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	ret := []string{}
	for {
		select {
		case <-doneCh:
			return ret, nil
		case str := <-outCh:
			ret = append(ret, str)
		case err := <-errCh:
			return nil, err
		}
	}
}
