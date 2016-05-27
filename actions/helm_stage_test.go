package actions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
	"github.com/arschles/sys"
	"github.com/deis/deisrel/testutil"
	"github.com/google/go-github/github"
)

func TestDownloadFiles(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "charts"
	var opt github.RepositoryContentGetOptions
	helmChart := WorkflowE2EChart
	helmChart.Files = []string{"README.md"}
	filePath := filepath.Join(helmChart.Name, helmChart.Files[0])

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, helmChart.Name), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `[
		  {
		    "type": "file",
		    "size": 625,
		    "name": "` + helmChart.Files[0] + `",
		    "path": "` + filePath + `",
				"content": "contents",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": "https://raw.githubusercontent.com/` + org + `/` + repo + `/master/` + filePath + `",
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  },
		  {
		    "type": "dir",
		    "size": 0,
		    "name": "` + helmChart.Name + `",
		    "path": "` + helmChart.Name + `",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": null,
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  }
		]`
		fmt.Fprintf(w, resp)
	})

	got, err := downloadFiles(ts.Client, org, repo, &opt, helmChart)

	assert.NoErr(t, err)
	assert.Equal(t, len(got), 1, "length of downloaded file slice")
	assert.Equal(t, got[0].FileName, helmChart.Files[0], "file name")
}

func TestDownloadFilesNotExist(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "charts"
	var opt github.RepositoryContentGetOptions
	helmChart := WorkflowE2EChart
	helmChart.Files = []string{"NotExist.md"}

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, helmChart.Name), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `[
		  {
		    "type": "dir",
		    "size": 0,
				"name": "` + helmChart.Name + `",
		    "path": "` + helmChart.Name + `",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": null,
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  }
		]`
		fmt.Fprintf(w, resp)
	})

	got, err := downloadFiles(ts.Client, org, repo, &opt, helmChart)

	assert.ExistsErr(t, err, "file doesn't exist")
	assert.Nil(t, got, "downloaded files slice")
}

func TestStageFiles(t *testing.T) {
	fakeFS := sys.NewFakeFS()

	readCloser := ioutil.NopCloser(bytes.NewBufferString(""))
	fileName := "testFile"
	ghFileToStage := ghFile{ReadCloser: readCloser, FileName: fileName}
	stagingDir := "staging"

	fakeFS.MkdirAll(stagingDir, os.ModePerm)
	stageFiles(fakeFS, []ghFile{ghFileToStage}, stagingDir)

	_, err := fakeFS.ReadFile(filepath.Join(stagingDir, fileName))
	assert.NoErr(t, err)
}

func TestCreateDir(t *testing.T) {
	fakeFS := sys.NewFakeFS()

	err := createDir(fakeFS, "foo")

	assert.NoErr(t, err)
	assert.Equal(t,
		fakeFS.Files,
		map[string]*bytes.Buffer{"foo": &bytes.Buffer{}},
		"file system contents")
}

func TestCreateDirAlreadyExists(t *testing.T) {
	fakeFS := sys.NewFakeFS()

	fakeFS.Create("foo")
	err := createDir(fakeFS, "foo")

	assert.NoErr(t, err)
	assert.Equal(t,
		fakeFS.Files,
		map[string]*bytes.Buffer{"foo": &bytes.Buffer{}},
		"file system contents")
}

func TestUpdateFilesWithRelease(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	fakeFP := sys.NewFakeFP()

	fileName := "foo/bar"
	fakeFS.Create(fileName)
	fakeFS.WriteFile(fileName, []byte("-dev"), os.ModePerm)
	var deisRelease = releaseName{
		Full:  "foobar",
		Short: "bar",
	}
	err := updateFilesWithRelease(fakeFP, fakeFS, deisRelease, fileName)

	assert.NoErr(t, err)
	actualFileContents, err := fakeFS.ReadFile(fileName)
	assert.NoErr(t, err)
	assert.Equal(t, actualFileContents, []byte("-"+deisRelease.Short), "updated file")
}

func TestUpdateFilesWithReleaseWithoutRelease(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	fakeFP := sys.NewFakeFP()

	fileName := "foo/bar"
	fakeFS.Create(fileName)
	fakeFS.WriteFile(fileName, []byte("-dev"), os.ModePerm)
	fakeFP.walkInvoked = false

	err := updateFilesWithRelease(fakeFP, fakeFS, deisRelease, fileName)
	assert.NoErr(t, err)
	assert.Equal(t, fakeFP.walkInvoked, false, "walk invoked")

	actualFileContents, err := fakeFS.ReadFile(fileName)
	assert.NoErr(t, err)
	assert.Equal(t, actualFileContents, []byte("-dev"), "updated file")
}
