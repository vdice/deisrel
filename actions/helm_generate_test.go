package actions

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

func TestGenerateParamsNoStage(t *testing.T) {
	fakeFS := getFakeFileSys()
	stagingDir := "staging/foo"
	defaultParamsComponentAttrs := genParamsComponentAttrs{
		Org:        "org",
		Tag:        "",
		PullPolicy: "stayGangsta",
	}
	paramsComponentMap := createParamsComponentMap()
	for _, componentName := range componentNames {
		paramsComponentMap[componentName] = defaultParamsComponentAttrs
	}

	err := generateParams(false, fakeFS, stagingDir, paramsComponentMap)
	assert.NoErr(t, err)

	expectedStagedFilepath := filepath.Join(stagingDir, "tpl/generate_params.toml")
	_, err = fakeFS.ReadFile(expectedStagedFilepath)
	assert.ExistsErr(t, err, "existence of staged file")
}

func TestGenerateParamsStage(t *testing.T) {
	fakeFS := getFakeFileSys()
	stagingDir := filepath.Join(stagingPath, "foo")
	defaultParamsComponentAttrs := genParamsComponentAttrs{
		Org:        "org",
		Tag:        "",
		PullPolicy: "stayGangsta",
	}
	paramsComponentMap := createParamsComponentMap()
	for _, componentName := range componentNames {
		paramsComponentMap[componentName] = defaultParamsComponentAttrs
	}

	err := generateParams(true, fakeFS, stagingDir, paramsComponentMap)
	assert.NoErr(t, err)

	expectedStagedFilepath := filepath.Join(stagingDir, "tpl/generate_params.toml")
	// verify file exists in fakeFS
	_, err = fakeFS.ReadFile(expectedStagedFilepath)
	assert.NoErr(t, err)

	actualFileContents, err := fakeFS.ReadFile(expectedStagedFilepath)
	assert.NoErr(t, err)
	expectedFileContents := new(bytes.Buffer)
	err = generateParamsTpl.Execute(expectedFileContents, paramsComponentMap)
	assert.NoErr(t, err)

	assert.Equal(t, actualFileContents, expectedFileContents.Bytes(), "staged file contents")

	// make sure each component name from canonical list exists in actualFileContents
	for _, componentName := range componentNames {
		lowerCasedComponentName := strings.ToLower(componentName)
		if lowerCasedComponentName != "workflowe2e" {
			if lowerCasedComponentName == "workflowmanager" {
				lowerCasedComponentName = "workflowManager"
			}
			assert.True(t,
				strings.Contains(string(actualFileContents), lowerCasedComponentName),
				fmt.Sprintf("component: %s not found!", lowerCasedComponentName))
		}
	}
}

func TestExecuteToStaging(t *testing.T) {
	fakeFS := getFakeFileSys()
	stagingDir := filepath.Join(stagingPath, "foo")

	_, err := executeToStaging(fakeFS, stagingDir)
	assert.NoErr(t, err)

	// just verify dir was created on fakeFS
	_, err = fakeFS.ReadFile(stagingDir)
	assert.NoErr(t, err)
}
