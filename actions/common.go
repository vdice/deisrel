package actions

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
		"workflow":         {"Workflow"},
		"workflow-e2e":     {"WorkflowE2E"},
		"workflow-manager": {"WorkflowManager"},
	}

	repoNames      = getRepoNames(repoToComponentNames)
	componentNames = getComponentNames(repoToComponentNames)
)

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
