package actions

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

type ghFile struct {
	ReadCloser io.ReadCloser
	FileName   string
}

var (
	ourFS = getRealFileSys()
	ourFP = getRealFilePath()
)

const (
	// RepoFlag represents the '-repo' flag
	RepoFlag = "repo"
	// RefFlag represents the '-ref' flag (for specifying a SHA, branch or tag)
	RefFlag = "ref"
	// GHOrgFlag represents the '-ghOrg' flag
	GHOrgFlag = "ghOrg"
)

func helmStage(ghClient *github.Client, c *cli.Context, fileNames []string, stagingSubDir string) {
	var opt github.RepositoryContentGetOptions
	opt.Ref = c.GlobalString(RefFlag)
	org := c.GlobalString(GHOrgFlag)
	repo := c.GlobalString(RepoFlag)

	ghFiles, err := downloadFiles(ghClient, org, repo, &opt, fileNames)
	if err != nil {
		log.Fatalf("Error downloading contents of %v (%s)", fileNames, err)
	}

	stageFiles(ourFS, ghFiles, stagingPath)

	if err := updateFilesWithRelease(ourFP, ourFS, deisRelease, stagingSubDir); err != nil {
		log.Fatalf("Error updating files with release '%s' (%s)", deisRelease.Short, err)
	}
}

func createDir(fs fileSys, dirName string) error {
	_, err := fs.Stat(dirName)
	if os.IsNotExist(err) {
		if err := fs.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func downloadFiles(ghClient *github.Client, org, repo string, opt *github.RepositoryContentGetOptions, fileNames []string) ([]ghFile, error) {
	ret := make([]ghFile, 0, len(fileNames))
	for _, fileName := range fileNames {
		readCloser, err := downloadContents(ghClient, org, repo, fileName, opt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ghFile{ReadCloser: readCloser, FileName: fileName})
	}
	return ret, nil
}

func stageFiles(fs fileSys, ghFiles []ghFile, stagingDir string) {
	for _, ghFile := range ghFiles {
		readCloser := ghFile.ReadCloser

		localFilePath := filepath.Join(stagingDir, ghFile.FileName)
		f, err := fs.Create(localFilePath)
		if err != nil {
			log.Fatalf("Error creating file %s (%s)", localFilePath, err)
		}

		if _, err := io.Copy(f, readCloser); err != nil {
			log.Fatalf("Error writing contents to file %s (%s)", localFilePath, err)
		}
		if err := f.Sync(); err != nil {
			log.Fatalf("Error flushing writes to stable storage (%s)", err)
		}
		log.Printf("File %s staged in '%s'", ghFile.FileName, stagingDir)
		defer readCloser.Close()
		defer f.Close()
	}
}

func updateFilesWithRelease(fp filePath, fs fileSys, release releaseName, walkPath string) error {
	if release.Full == "" || release.Short == "" {
		log.Printf("DEIS_RELEASE (%s) and/or DEIS_RELEASE_SHORT (%s) not provided so not amending staged files.", release.Full, release.Short)
	} else {
		if err := fp.Walk(walkPath, getReleaseWalker().handlerFunc(fs, release)); err != nil {
			return err
		}
	}
	return nil
}
