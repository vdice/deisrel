package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/testutil"
)

func TestCheckExistenceOnHubFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	testImageAndTag := ImageAndTag{
		Image: "myImage",
		Tag:   "myTag",
	}

	ts.Mux.HandleFunc(fmt.Sprintf("/v2/%s/tags/list", testImageAndTag.Image), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `{
			"name": "` + testImageAndTag.Image + `",
			"tags": [
				"` + testImageAndTag.Tag + `"
			]
		}`
		fmt.Fprintf(w, resp)
	})

	fakeHubClient := NewFakeHubRegistry(t, ts)
	err := fakeHubClient.CheckExistence(testImageAndTag)
	assert.NoErr(t, err)
}

func TestCheckExistenceOnHubNotFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	testImageAndTag := ImageAndTag{
		Image: "myImage",
		Tag:   "myTag",
	}

	ts.Mux.HandleFunc(fmt.Sprintf("/v2/%s/tags/list", testImageAndTag.Image), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `{
			"name": "` + testImageAndTag.Image + `",
			"tags": [
				"notMyTag"
			]
		}`
		fmt.Fprintf(w, resp)
	})

	fakeHubClient := NewFakeHubRegistry(t, ts)
	err := fakeHubClient.CheckExistence(testImageAndTag)

	expectedErr := ErrTagNotFound{
		imageAndTag: testImageAndTag,
		registry:    fakeHubClient.Client.URL,
	}
	assert.Err(t, expectedErr, err)
}

func TestPushTagToHub(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	arch := "arch"
	blobSum := "sha256:1234abcd"
	history := "some-history"
	orig := ImageAndTag{
		Image: "origImage",
		Tag:   "origTag",
	}
	new := ImageAndTag{
		Image: "newImage",
		Tag:   "newTag",
	}

	// need handler for get (orig) manifest
	ts.Mux.HandleFunc(fmt.Sprintf("/v2/%s/manifests/%s", orig.Image, orig.Tag), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}

		resp := `{
		   "schemaVersion": 1,
		   "name": "` + orig.Image + `",
		   "tag": "` + orig.Tag + `",
		   "architecture": "` + arch + `",
		   "fsLayers": [
		      {
		         "blobSum": "` + blobSum + `"
		      }
		   ],
		   "history": [
		      {
		         "v1Compatibility": "` + history + `"
		      }
		   ],
		   "signatures": []
		}`
		fmt.Fprintf(w, resp)
	})

	// and handler for put (new) manifest
	ts.Mux.HandleFunc(fmt.Sprintf("/v2/%s/manifests/%s", new.Image, new.Tag), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "PUT" {
			t.Errorf("Request method: %v, want PUT", got)
		}

		defer r.Body.Close()
		contents, err := ioutil.ReadAll(r.Body)
		assert.NoErr(t, err)

		type Got struct {
			Name     string              `json:"name"`
			Tag      string              `json:"tag"`
			Arch     string              `json:"architecture"`
			FSLayers []map[string]string `json:"fsLayers"`
			History  []map[string]string `json:"history"`
		}

		g := Got{}
		err = json.Unmarshal(contents, &g)
		assert.NoErr(t, err)

		assert.Equal(t, g.Name, new.Image, "image name")
		assert.Equal(t, g.Tag, new.Tag, "tag value")
		assert.Equal(t, g.Arch, arch, "architecture value")
		assert.Equal(t, g.FSLayers[0]["blobSum"], blobSum, "blobSum value")
		assert.Equal(t, g.History[0]["v1Compatibility"], history, "history value")
	})

	fakeHubClient := NewFakeHubRegistry(t, ts)
	err := fakeHubClient.PushTag(orig, new)
	assert.NoErr(t, err)
}
