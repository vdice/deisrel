package actions

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	sys "github.com/arschles/sys"
)

//TODO: PR these additions into arschles/sys and remove depending on what gets in
type fileSys interface {
	sys.FS
	Create(string) (io.WriteCloser, error)
	Stat(string) (os.FileInfo, error)
	MkdirAll(string, os.FileMode) error
	WriteFile(string, []byte, os.FileMode) (int, error)
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

func (r *realFileSys) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (r *realFileSys) WriteFile(filename string, data []byte, perm os.FileMode) (int, error) {
	return len(data), ioutil.WriteFile(filename, data, perm)
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
	Files map[string]*bytes.Buffer
}

func getFakeFileSys() *fakeFileSys {
	return &fakeFileSys{Files: make(map[string]*bytes.Buffer)}
}

func (f *fakeFileSys) ReadFile(name string) ([]byte, error) {
	buf, ok := f.Files[name]
	if !ok {
		return nil, sys.FakeFileNotFound{Filename: name}
	}
	return buf.Bytes(), nil
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

type inMemoryCloser struct {
	buf *bytes.Buffer
}

func (i inMemoryCloser) Write(b []byte) (int, error) {
	return i.buf.Write(b)
}

func (i inMemoryCloser) Close() error {
	return nil
}

func (f *fakeFileSys) Create(path string) (io.WriteCloser, error) {
	buf := new(bytes.Buffer)
	f.Files[path] = buf
	return inMemoryCloser{buf: buf}, nil
}

func (f *fakeFileSys) WriteFile(filename string, data []byte, perm os.FileMode) (int, error) {
	// clear out old contents as Buffer.Write appends
	f.Files[filename] = new(bytes.Buffer)
	return f.Files[filename].Write(data)
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
