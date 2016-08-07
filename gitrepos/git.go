package gitrepos

import (
	"io"
	"os"
	"os/exec"
	"strings"
	_path "path"
)

type Git struct {
	Repo string
	BinPath string
}

func NewGit(repo string) *Git {
	return &Git{
		Repo: repo,
		BinPath: "/usr/bin/git",
	}
}

func (g *Git) Command(arg ...string) *exec.Cmd {
	path := "PATH=" + os.Getenv("PATH")

	return &exec.Cmd{
		Path: g.BinPath,
		Args: append([]string{g.BinPath}, arg...),
		Dir: g.Repo,
		Env: []string{path},
	}
}

func (g *Git) ListFiles(path string) ([]os.FileInfo, error) {
	output, err := g.Command("ls-tree", "HEAD", path).Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	files := make([]os.FileInfo, 0, len(lines))

	for _, line := range(lines) {
		cols := splitRegexp.Split(line, 4)

		if len(cols) != 4 {
			continue
		}

		fi := &FileInfo{
			mode: cols[0],
			objectType: cols[1],
			hash: cols[2],
			name: _path.Base(cols[3]),
		}

		files = append(files, fi)
	}

	return files, nil
}

func (g *Git) OpenFile(hash string) (*exec.Cmd, io.Reader, error) {
	cmd := g.Command("cat-file", "blob", hash)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()

	if err != nil {
		return nil, nil, err
	}

	return cmd, stdout, nil
}

func (g *Git) IsGitRepo() (bool, error) {
	err := g.Command("show").Run()

	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return false, err
		}
	}

	return err == nil, nil
}

