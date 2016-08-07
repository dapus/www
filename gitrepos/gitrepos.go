package gitrepos

import (
	"net/http"
	"os"
	"strings"
	"regexp"
	"errors"
)

var splitRegexp *regexp.Regexp

func init() {
	splitRegexp = regexp.MustCompile("[ \t]")
}

// Implements http.FileServer interface
type GitRepos string

func (d GitRepos) openRoot() (http.File, error) {
	f, err := os.Open(string(d))
	defer f.Close()

	if err != nil {
		return nil, err
	}

	fileInfo, err := f.Stat()

	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	{
		f, err := f.Readdir(0)

		if err != nil {
			return nil, err
		}

		for _, f := range(f) {
			if !f.IsDir() {
				continue
			}

			path := []string{string(d), f.Name()}
			git := NewGit(strings.Join(path, "/"))

			if ok, err := git.IsGitRepo(); err != nil {
				return nil, err
			} else if ok {
				files = append(files, f)
			}
		}
	}

	file := &File{
		name: "/",
		file: fileInfo,
		files: files,
	}

	return file, nil
}

func (d GitRepos) Open(name string) (http.File, error) {
	if len(name) > 0 && name[len(name)-1:] == "/" {
		name = name[:len(name)-1]
	}

	pathSeg := strings.Split(name, "/")

	if len(pathSeg) < 2 {
		return d.openRoot()
	}

	repo := strings.Join(pathSeg[:2], "/")
	name = strings.Join(pathSeg[2:], "/")

	file := &File{
		git: NewGit(string(d) + "/" + repo),
		name: name,
	}

	files, err := file.git.ListFiles(name)

	if err != nil {
		return nil, err
	}

	if len(files) < 1 {
		return nil, errors.New("no such file")
	}

	if name != "" && len(files) > 1 {
		panic("Number of files are greater than one")
	}

	if name == "" {
		file.file = &FileInfo{
			objectType: "tree",
			name: "/",
		}
	} else {
		file.file = files[0]
	}

	if file.file.IsDir() {
		path := name

		if name == "" {
			path = "."
		}

		files, err := file.git.ListFiles(path + "/")

		if err != nil {
			return nil, err
		}

		file.files = files

		return file, nil
	}

	cmd, r, err := file.git.OpenFile(file.file.(*FileInfo).hash)

	if err != nil {
		return nil, err
	}

	file.cmd = cmd
	file.out = r

	return file, nil
}
