package registry

import (
	"log"
	"os"
	"strings"

	"github.com/coreos/go-quay/models"
	"github.com/coreos/go-quay/quay"
	"github.com/coreos/go-quay/quay/tag"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const (
	quayTokenEnvVarName = "QUAY_AUTH_TOKEN"
)

// QuayRegistry includes a Client, Auth and helper methods for interacting
// with the quay.io api
type QuayRegistry struct {
	Client *quay.Client
	Auth   QuayAuth
}

// NewQuayRegistry returns a pointer to a new QuayRegistry
func NewQuayRegistry(c *quay.Client, auth QuayAuth) *QuayRegistry {
	return &QuayRegistry{
		Client: c,
		Auth:   auth,
	}
}

// QuayAuth is an interface providing AuthenticateRequest
type QuayAuth interface {
	AuthenticateRequest(runtime.ClientRequest, strfmt.Registry) error
}

// NewQuayAuth returns a QuayAuth instance
func NewQuayAuth() QuayAuth {
	return &realQuayAuth{}
}

type realQuayAuth struct{}

func (rqa *realQuayAuth) AuthenticateRequest(cr runtime.ClientRequest, reg strfmt.Registry) error {
	quayTkn := os.Getenv(quayTokenEnvVarName)
	if quayTkn == "" {
		log.Fatalf("'%s' env var required", quayTokenEnvVarName)
	}
	return cr.SetHeaderParam("Authorization", quayTkn)
}

// CheckExistence is the ExistenceChecker for QuayRegistry
func (q *QuayRegistry) CheckExistence(imageAndTag ImageAndTag) error {
	quayListTagImagesParams := &tag.ListTagImagesParams{
		Repository: imageAndTag.Image,
		Tag:        imageAndTag.Tag,
	}

	_, err := q.Client.Tag.ListTagImages(quayListTagImagesParams, q.Auth)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return ErrTagNotFound{
				imageAndTag: imageAndTag,
				registry:    "quay.io",
			}
		} else if strings.Contains(err.Error(), "write on closed buffer") {
			// TODO: why are we getting these?  Safe to drop?
			// log.Printf("Encountered 'write on closed buffer' for %s", imageAndTag.Image)
			return nil
		}
	}
	return err
}

// PushTag is the TagPusher for QuayRegistry
func (q *QuayRegistry) PushTag(orig ImageAndTag, new ImageAndTag) error {
	origFullName := orig.GetFullName()
	changeTagImageParams := &tag.ChangeTagImageParams{
		Body:       &models.MoveTag{Image: &origFullName},
		Repository: new.Image,
		Tag:        new.Tag,
	}

	_, err := q.Client.Tag.ChangeTagImage(changeTagImageParams, q.Auth)
	// if err != nil {
	// 	return err
	// }
	return err
}
