package actions

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error {
	return nil
}

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func generateParams(stage bool, fs fileSys, whereTo string, paramsComponentMap genParamsComponentMap, template *template.Template) error {
	var executeTo = NopWriteCloser(os.Stdout)
	if stage {
		defer executeTo.Close()
		var err error
		executeTo, err = executeToStaging(fs, filepath.Join(whereTo, "tpl"))
		if err != nil {
			log.Fatalf("Error creating staging file (%s)", err)
		}
	}
	return template.Execute(executeTo, paramsComponentMap)
}

func executeToStaging(fs fileSys, stagingSubDir string) (io.WriteCloser, error) {
	fs.MkdirAll(stagingSubDir, os.ModePerm)
	return fs.Create(filepath.Join(stagingSubDir, generateParamsFileName))
}
