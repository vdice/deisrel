package registry

import (
	"fmt"
	"strings"
)

// Registry represents a general registry interface, currently mandating that
// each implementation define CheckExistence functionality
type Registry interface {
	CheckExistence(imgTag ImageAndTag) error
	PushTag(orig ImageAndTag, new ImageAndTag) error
}

// ExistenceChecker checks for existence of imgTag and returns error
type ExistenceChecker func(imgTag ImageAndTag) error

// TagPusher pushes new ImageAndTag based on orig ImageAndTag and returns error
type TagPusher func(orig ImageAndTag, new ImageAndTag) error

// ImageAndTag represents a Docker Image and Tag pair
type ImageAndTag struct {
	Image string
	Tag   string
}

// GetFullName returns the full <it.Image>:<it.Tag> string representation
func (it *ImageAndTag) GetFullName() string {
	return strings.Join([]string{it.Image, it.Tag}, ":")
}

// ErrTagNotFound is the error used when an imageAndTag pair is not found in a
// given registry
type ErrTagNotFound struct {
	imageAndTag ImageAndTag
	registry    string
}

func (e ErrTagNotFound) Error() string {
	return fmt.Sprintf("Tag '%s' not found for image '%s' from registry %s",
		e.imageAndTag.Tag, e.imageAndTag.Image, e.registry)
}
