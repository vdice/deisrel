package actions

import (
	"text/template"
)

const (
	// TODO: https://github.com/deis/deisrel/issues/11
	generateParamsE2ETplStr = `[e2e]
org = "{{.WorkflowE2E.Org}}"
dockerTag = "{{.WorkflowE2E.Tag}}"
pullPolicy = "{{.WorkflowE2E.PullPolicy}}"
`
)

var (
	generateParamsE2ETpl = template.Must(template.New("generateParamsE2ETpl").Parse(generateParamsE2ETplStr))
)
