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
	"github.com/deis/deisrel/testutil"
	"github.com/google/go-github/github"
)

func TestDownloadFiles(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "chart"
	var opt github.RepositoryContentGetOptions
	fileDir := "workflow-dev-e2e"
	fileName := "README.md"
	filePath := filepath.Join(fileDir, fileName)

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, fileDir), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `[
		  {
		    "type": "file",
		    "size": 625,
		    "name": "` + fileName + `",
		    "path": "` + filePath + `",
				"content": "contents",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": "https://raw.githubusercontent.com/` + org + `/` + repo + `master/` + filePath + `",
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  },
		  {
		    "type": "dir",
		    "size": 0,
		    "name": "` + fileDir + `",
		    "path": "` + fileDir + `",
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

	got, err := downloadFiles(ts.Client, org, repo, &opt, []string{filePath})

	assert.NoErr(t, err)
	assert.Equal(t, len(got), 1, "length of downloaded file slice")
	assert.Equal(t, got[0].FileName, filePath, "file name")
}

func TestDownloadFilesNotExist(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "chart"
	var opt github.RepositoryContentGetOptions
	fileDir := "workflow-dev-e2e"
	fileName := "NotExist.md"
	filePath := filepath.Join(fileDir, fileName)

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, fileDir), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `[
		  {
		    "type": "dir",
		    "size": 0,
		    "name": "` + fileDir + `",
		    "path": "` + fileDir + `",
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

	got, err := downloadFiles(ts.Client, org, repo, &opt, []string{filePath})

	assert.ExistsErr(t, err, "file doesn't exist")
	assert.Nil(t, got, "downloaded files slice")
}

func TestStageFiles(t *testing.T) {
	fakeFileSys := getFakeFileSys()

	readCloser := ioutil.NopCloser(bytes.NewBufferString(""))
	fileName := "testFile"
	ghFileToStage := ghFile{ReadCloser: readCloser, FileName: fileName}
	stagingDir := "staging"

	fakeFileSys.MkdirAll(stagingDir, os.ModePerm)
	stageFiles(fakeFileSys, []ghFile{ghFileToStage}, stagingDir)

	assert.Equal(t,
		fakeFileSys.Files,
		map[string][]byte{
			stagingDir:                          []byte{},
			filepath.Join(stagingDir, fileName): []byte{}},
		"file system contents")
}

func TestCreateDir(t *testing.T) {
	fakeFileSys := getFakeFileSys()

	err := createDir(fakeFileSys, "foo")

	assert.NoErr(t, err)
	assert.Equal(t,
		fakeFileSys.Files,
		map[string][]byte{"foo": []byte{}},
		"file system contents")
}

func TestCreateDirAlreadyExists(t *testing.T) {
	fakeFileSys := getFakeFileSys()

	fakeFileSys.Files["foo"] = []byte{}
	err := createDir(fakeFileSys, "foo")

	assert.NoErr(t, err)
	assert.Equal(t,
		fakeFileSys.Files,
		map[string][]byte{"foo": []byte{}},
		"file system contents")
}

func TestUpdateFilesWithRelease(t *testing.T) {
	fakeFileSys := getFakeFileSys()
	fakeFilePath := getFakeFilePath()

	fileName := "foo/bar"
	fakeFileSys.Create(fileName)
	fakeFileSys.WriteFile(fileName, []byte("dev"), os.ModePerm)
	var deisRelease = releaseName{
		Full:  "foobar",
		Short: "bar",
	}
	err := updateFilesWithRelease(fakeFilePath, fakeFileSys, deisRelease, fileName)

	assert.NoErr(t, err)
	assert.Equal(t, fakeFileSys.Files[fileName], []byte(deisRelease.Short), "updated file")
}

func TestUpdateFilesWithReleaseWithoutRelease(t *testing.T) {
	fakeFileSys := getFakeFileSys()
	fakeFilePath := getFakeFilePath()

	fileName := "foo/bar"
	fakeFileSys.Create(fileName)
	fakeFileSys.WriteFile(fileName, []byte("dev"), os.ModePerm)
	fakeFilePath.walkInvoked = false
	err := updateFilesWithRelease(fakeFilePath, fakeFileSys, deisRelease, fileName)

	assert.NoErr(t, err)
	assert.Equal(t, fakeFilePath.walkInvoked, false, "walk invoked")
	assert.Equal(t, fakeFileSys.Files[fileName], []byte("dev"), "updated file")
}
