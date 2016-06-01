package actions

import (
	"text/template"
)

const (
	// TODO: https://github.com/deis/deisrel/issues/11
	generateParamsRouterTplStr = `[router]
org = "{{.Router.Org}}"
pullPolicy = "{{.Router.PullPolicy}}"
dockerTag = "{{.Router.Tag}}"
platformDomain = ""
`
)

var (
	generateParamsRouterTpl = template.Must(template.New("generateParamsRouterTpl").Parse(generateParamsRouterTplStr))
)
