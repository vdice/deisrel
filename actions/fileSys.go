package actions

import (
	"io/ioutil"
	"os"
	"path/filepath"

	sys "github.com/arschles/sys"
	"github.com/deis/deisrel/testutil"
)

//TODO: PR these additions into arschles/sys and remove depending on what gets in
type fileSys interface {
	sys.FS
	Create(string) (*os.File, error)
	Stat(string) (os.FileInfo, error)
	MkdirAll(string, os.FileMode) error
	WriteFile(string, []byte, os.FileMode) error
}

func getRealFileSys() fileSys {
	return &realFileSys{}
}

type realFileSys struct{}

func (r *realFileSys) ReadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func (r *realFileSys) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

func (r *realFileSys) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (r *realFileSys) MkdirAll(dirName string, perm os.FileMode) error {
	return os.MkdirAll(dirName, perm)
}

func (r *realFileSys) Create(path string) (*os.File, error) {
	return os.Create(path)
}

func (r *realFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

type fakeFileInfo struct {
	os.FileInfo
	isDir bool
}

func getFakeFileInfo() *fakeFileInfo {
	return &fakeFileInfo{}
}

func (ffi *fakeFileInfo) IsDir() bool {
	return ffi.isDir
}

type fakeFileSys struct {
	sys.FakeFS
	Files map[string][]byte
}

func getFakeFileSys() *fakeFileSys {
	return &fakeFileSys{Files: make(map[string][]byte)}
}

func (f *fakeFileSys) ReadFile(name string) ([]byte, error) {
	b, ok := f.Files[name]
	if !ok {
		return nil, sys.FakeFileNotFound{Filename: name}
	}
	return b, nil
}
func (f *fakeFileSys) Stat(path string) (os.FileInfo, error) {
	_, err := f.ReadFile(path)
	if err != nil {
		return nil, os.ErrNotExist
	}
	return getFakeFileInfo(), nil
}

func (f *fakeFileSys) MkdirAll(dirName string, perm os.FileMode) error {
	_, err := f.Create(dirName)
	return err
}

func (f *fakeFileSys) Create(path string) (*os.File, error) {
	f.Files[path] = []byte{}
	ret, err := os.OpenFile(filepath.Join(testutil.TestDataDir(), "foo.txt"), os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (f *fakeFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if _, err := f.Stat(filename); err != nil {
		return err
	}
	f.Files[filename] = data
	return nil
}

type filePath interface {
	Walk(string, filepath.WalkFunc) error
}

func getRealFilePath() filePath {
	return &realFilePath{}
}

type realFilePath struct {
	filePath
}

func (r *realFilePath) Walk(root string, walkFunc filepath.WalkFunc) error {
	return filepath.Walk(root, walkFunc)
}

type fakeFilePath struct {
	fakeFileSys
	walkInvoked bool
}

func getFakeFilePath() *fakeFilePath {
	return &fakeFilePath{}
}

func (f *fakeFilePath) Walk(root string, walkFunc filepath.WalkFunc) error {
	f.walkInvoked = true
	fi := getFakeFileInfo()
	fi.isDir = false
	return walkFunc(root, fi, nil)
}
