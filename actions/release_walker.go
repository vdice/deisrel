package actions

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/arschles/sys"
)

type releaseWalker struct {
	filepath.WalkFunc
	sys.FakeFS
}

// RegexReplacer is a struct containing a *regexp.Regexp representing the string
// to replace, along with repl, the replacement string.  Modeled after strings/Replacer.
type RegexReplacer struct {
	regex *regexp.Regexp
	repl  string
}

// NewRegexReplacer returns a new RegexReplacer given a src string to convert
// to a *regexp.Regexp instance and a repl string
func NewRegexReplacer(src string, repl string) RegexReplacer {
	return RegexReplacer{
		regex: regexp.MustCompile(src),
		repl:  repl,
	}
}

func getReleaseWalker() *releaseWalker {
	return &releaseWalker{}
}

func (r *releaseWalker) handlerFunc(fs sys.FS, release releaseName) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		return r.walk(path, release, fi, err, fs)
	}
}

func (r *releaseWalker) walk(path string, release releaseName, fi os.FileInfo, err error, fs sys.FS) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	read, err := fs.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file %s (%s)", path, err)
	}

	regexReplacers := []RegexReplacer{
		NewRegexReplacer("-dev", "-"+release.Short),
		NewRegexReplacer("v[0-9].[0-9].[0-9]", release.Full),
		NewRegexReplacer("WARNING: this chart is for testing only! Features may not work and there are likely to be bugs.\n", ""),
		NewRegexReplacer("\\s\\(For testing only!\\)", ""),
	}

	if !(release.Full == release.Short) {
		regexReplacers = append(regexReplacers, NewRegexReplacer("v[0-9].[0-9].[0-9]-dev", release.Full+"-"+release.Short))
	}

	newContents := string(read)
	for _, regexReplacer := range regexReplacers {
		newContents = regexReplacer.regex.ReplaceAllString(newContents, regexReplacer.repl)
	}

	if _, err := fs.WriteFile(path, []byte(newContents), 0); err != nil {
		log.Fatalf("Error writing contents to file %s (%s)", path, err)
	}

	log.Printf("File '%s' updated with release '%s'", path, release.Short)
	return nil
}
