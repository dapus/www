package gitrepos

import (
	"os"
	"time"
)

type FileInfo struct {
	mode string
	objectType string
	hash string
	name string
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return -1
}

func (f *FileInfo) Mode() os.FileMode {
	return os.FileMode(0)
}

func (f *FileInfo) ModTime() time.Time {
	return time.Now()
}

func (f *FileInfo) IsDir() bool {
	return f.objectType == "tree"
}

func (f *FileInfo) Sys() interface{} {
	return nil
}
