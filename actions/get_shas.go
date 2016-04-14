package actions

import (
	"fmt"
	"log"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// GetShas is the CLI action for getting github shas of all of the Deis Workflow repos
func GetShas(c *cli.Context) {
	ghClient := github.NewClient(nil)
	outCh := make(chan string)
	var wg sync.WaitGroup
	for _, repo := range repoNames {
		wg.Add(1)
		ch := make(chan string)
		go func(repo string) {
			defer wg.Done()
			repoCommits, _, err := ghClient.Repositories.ListCommits("deis", repo, &github.CommitsListOptions{
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			})
			if err != nil {
				log.Fatalf("Error listing commits for repo %s (%s)", repo, err)
			}
			if len(repoCommits) < 1 {
				log.Fatalf("No commits found for repo %s", repo)
			}
			repoCommit := repoCommits[0]
			sha := *repoCommit.SHA
			ch <- fmt.Sprintf("%s: %s\n", repo, sha[0:7])
		}(repo)
		go func() {
			outCh <- <-ch
		}()
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()

	for str := range outCh {
		fmt.Print(str)
	}
}
