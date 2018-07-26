package commands

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path"
	"syscall"

	"github.com/ericaro/frontmatter"
)

// CommandRunner is an interface for runnable commands
type CommandRunner interface {
	Run(context.Context, []string)
}

// CommandContextKey is a context key specific to commands package
type CommandContextKey string

func gitCommand(params ...string) int {
	cmd := exec.Command("git", params...)
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
	}
	return 0
}

type entryHeader struct {
	Filepath string   `yaml:"-"`
	Filename string   `yaml:"-"`
	Tags     []string `yaml:"tags"`
	Content  string   `fm:"content" yaml:"-"`
}

type frontmatterResult struct {
	header *entryHeader
	err    error
}

func readFrontmatter(filePath string, results chan<- frontmatterResult) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		results <- frontmatterResult{
			header: nil,
			err:    err,
		}
		return
	}
	head := new(entryHeader)
	frontmatter.Unmarshal(content, head)
	head.Filepath = filePath
	head.Filename = path.Base(filePath)
	results <- frontmatterResult{
		header: head,
		err:    nil,
	}
}
