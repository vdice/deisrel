package registry

import hub "github.com/heroku/docker-registry-client/registry"

// HubRegistry consists of a client and helper methods for interacting
// with the DockerHub api
type HubRegistry struct {
	Client *hub.Registry
}

// NewHubRegistry returns a pointer to a new HubRegistry
func NewHubRegistry(reg *hub.Registry) *HubRegistry {
	return &HubRegistry{Client: reg}
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

// GetHub returns a new *hub.Registry using the provided parameters
func GetHub(url, username, password string) (*hub.Registry, error) {
	return hub.New(url, username, password)
}

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
