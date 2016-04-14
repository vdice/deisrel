package actions

import (
	"text/template"
)

const (
	generateParamsE2ETplStr = `[e2e]
org = "{{.Org}}"
dockerTag = "{{.Tag}}"
pullPolicy = "{{.PullPolicy}}"`
)

var (
	generateParamsE2ETpl = template.Must(template.New("generateParamsE2ETpl").Parse(generateParamsE2ETplStr))
)
