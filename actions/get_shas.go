package actions

import (
	"fmt"
	"log"
	"sync"

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
				if c.Bool(ShortFlag) {
					sha = sha[:7]
				}
				ch <- fmt.Sprintf("%s: %s\n", repo, sha)
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
}
