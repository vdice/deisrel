package testutil

import (
	"log"
	"os"
)

var goPath string

func init() {
	goPath = os.Getenv("GOPATH")
	if goPath == "" {
		log.Fatalf("GOPATH not set")
	}
}
