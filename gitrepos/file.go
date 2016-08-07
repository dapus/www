package gitrepos

import (
	"io"
	"os"
	"os/exec"
	"errors"
)

type File struct {
	git *Git
	files []os.FileInfo
	file os.FileInfo
	d string
	name string
	cmd *exec.Cmd
	out io.Reader
}

func (f *File) Close() error {
	if f.cmd == nil {
		return nil
	}

	err := f.cmd.Wait()
	f.cmd = nil

	return err
}

func (f *File) Read(p []byte) (int, error) {
	if f.file.IsDir() {
		return 0, errors.New("not a regularfiles")
	}

	return f.out.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return offset, nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if !f.file.IsDir() {
		return nil, errors.New("not a directory")
	}

	var files = f.files

	if files == nil {
		var err error

		files, err = f.git.ListFiles(f.name)

		if err != nil {
			return nil, err
		}
	}

	if count > 0 {
		return files[0:count], nil
	}

	return files, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.file, nil;
}

