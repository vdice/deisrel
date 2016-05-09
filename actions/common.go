package actions

import (
	"log"
	"os"
	"path/filepath"
)

const (
	// TagFlag represents the '-tag' flag
	TagFlag = "tag"
	// PullPolicyFlag represents the '-pull-policy' flag
	PullPolicyFlag = "pull-policy"
	// OrgFlag represents the '-org' flag
	OrgFlag = "org"
	// ShaFilepathFlag represents the --sha-filepath flag
	ShaFilepathFlag = "sha-filepath"
	// YesFlag represents the --yes flag
	YesFlag = "yes"
	// StageFlag represents the '-stage' flag
	StageFlag = "stage"
)

const (
	generateParamsFileName = "generate_params.toml"
)

var (
	// TODO: https://github.com/deis/deisrel/issues/12
	repoToComponentNames = map[string][]string{
		"builder":          {"Builder"},
		"controller":       {"Controller"},
		"dockerbuilder":    {"DockerBuilder"},
		"fluentd":          {"FluentD"},
		"monitor":          {"InfluxDB", "Grafana", "Telegraf"},
		"logger":           {"Logger"},
		"minio":            {"Minio"},
		"postgres":         {"Database"},
		"registry":         {"Registry"},
		"router":           {"Router"},
		"slugbuilder":      {"SlugBuilder"},
		"slugrunner":       {"SlugRunner"},
		"stdout-metrics":   {"StdoutMetrics"},
		"workflow":         {"Workflow"},
		"workflow-e2e":     {"WorkflowE2E"},
		"workflow-manager": {"WorkflowManager"},
	}

	repoNames      = getRepoNames(repoToComponentNames)
	componentNames = getComponentNames(repoToComponentNames)
	deisRelease    = releaseName{
		Full:  os.Getenv("DEIS_RELEASE"),
		Short: os.Getenv("DEIS_RELEASE_SHORT"),
	}
	stagingPath = getFullPath("staging")
)

type releaseName struct {
	Full  string
	Short string
}

func getRepoNames(repoToComponentNames map[string][]string) []string {
	repoNames := make([]string, 0, len(repoToComponentNames))
	for repoName := range repoToComponentNames {
		repoNames = append(repoNames, repoName)
	}
	return repoNames
}

func getComponentNames(repoToComponentNames map[string][]string) []string {
	var ret []string
	for _, componentNames := range repoToComponentNames {
		for _, componentName := range componentNames {
			ret = append(ret, componentName)
		}
	}
	return ret
}

func getFullPath(dirName string) string {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working dir (%s)", err)
	}
	return filepath.Join(currentWorkingDir, dirName)
}
