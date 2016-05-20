package registry

import (
	"github.com/docker/distribution/manifest"
	"github.com/docker/libtrust"
	hub "github.com/heroku/docker-registry-client/registry"
)

// HubRegistry consists of a client and helper methods for interacting
// with the DockerHub api
type HubRegistry struct {
	Client *hub.Registry
}

// NewHubRegistry returns a pointer to a new HubRegistry
func NewHubRegistry(reg *hub.Registry) *HubRegistry {
	return &HubRegistry{Client: reg}
}

// GetHub returns a new *hub.Registry client using the provided parameters
func GetHub(url, username, password string) (*hub.Registry, error) {
	return hub.New(url, username, password)
}

// CheckExistence is the ExistenceChecker for HubRegistry
func (h *HubRegistry) CheckExistence(imageAndTag ImageAndTag) error {
	tagsFromHub, err := h.Client.Tags(imageAndTag.Image)
	if err != nil {
		return err
	}

	if !sliceContains(tagsFromHub, imageAndTag.Tag) {
		return ErrTagNotFound{
			imageAndTag: imageAndTag,
			registry:    h.Client.URL,
		}
	}
	return nil
}

// PushTag is the TagPusher for HubRegistry
func (h *HubRegistry) PushTag(orig ImageAndTag, new ImageAndTag) error {
	origManifest, err := h.Client.Manifest(orig.Image, orig.Tag)
	if err != nil {
		return err
	}
	// fmt.Println(origManifest)

	newManifest := &manifest.Manifest{
		Versioned: manifest.Versioned{
			SchemaVersion: 1,
		},
		Name:         new.Image,
		Tag:          new.Tag,
		FSLayers:     origManifest.FSLayers,
		Architecture: origManifest.Architecture,
		History:      origManifest.History,
	}

	key, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return err
	}

	signedManifest, err := manifest.Sign(newManifest, key)
	if err != nil {
		return err
	}

	// fmt.Println(signedManifest)
	if err := h.Client.PutManifest(new.Image, new.Tag, signedManifest); err != nil {
		return err
	}

	return nil
}
