package actions

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"
)

type repoAndSha struct {
	repoName string
	sha      string
}

func noTransform(s string) string       { return s }
func shortShaTransform(s string) string { return s[:7] }

func getShas(ghClient *github.Client, repos []string, transform func(string) string) ([]repoAndSha, error) {
	outCh := make(chan repoAndSha)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		ch := make(chan repoAndSha)
		ech := make(chan error)
		go func(repo string) {
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
			ch <- repoAndSha{repoName: repo, sha: sha}
		}(repo)
		go func() {
			defer wg.Done()
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

	ret := []repoAndSha{}
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

func getLastTag(ghClient *github.Client, repos []string) (map[string]*github.RepositoryTag, error) {
	for _, repo := range repos {
		_, _, err := ghClient.Repositories.ListTags("deis", repo, nil)
		if err != nil {
			return make(map[string]*github.RepositoryTag), err
		}
	}
	return make(map[string]*github.RepositoryTag), nil
}
