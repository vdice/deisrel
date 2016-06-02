package actions

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/changelog"
	"github.com/deis/deisrel/testutil"
)

func TestGenerateChangelogGlobal(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	got := &changelog.Values{
		OldRelease: "old",
		NewRelease: "new",
	}

	for _, repoName := range allGitRepoNames {
		ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/compare/old...new", repoName), func(w http.ResponseWriter, r *http.Request) {
			if got := r.Method; got != "GET" {
				t.Errorf("Request method: %v, want GET", got)
			}
			fmt.Fprintf(w, `{
			  "base_commit": {
			    "sha": "s",
			    "commit": {
			      "author": { "name": "n" },
			      "committer": { "name": "n" },
			      "message": "m",
			      "tree": { "sha": "t" }
			    },
			    "author": { "login": "n" },
			    "committer": { "login": "l" },
			    "parents": [ { "sha": "s" } ]
			  },
			  "status": "s",
			  "ahead_by": 1,
			  "behind_by": 2,
			  "total_commits": 1,
			  "commits": [
			    {
			      "sha": "abc1234567890",
			      "commit": { "author": { "name": "n" }, "message": "feat(foo): new feature!" },
			      "author": { "login": "l" },
			      "committer": { "login": "l" },
			      "parents": [ { "sha": "s" } ]
			    }
			  ],
			  "files": [ { "filename": "f" } ]
			}`)
		})
	}

	changelogVals, errs := generateChangelogVals(ts.Client, got.OldRelease, got.NewRelease)
	for _, err := range errs {
		assert.Nil(t, err, "generateChangelogVals err")
	}
	assert.Equal(t, len(changelogVals), len(allGitRepoNames), "number of repos checked")
}
