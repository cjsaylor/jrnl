package commands

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/ericaro/frontmatter"
)

// CommandRunner is an interface for runnable commands
type CommandRunner interface {
	Run(context.Context, []string) error
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

const JournalTimeformat = "Mon Jan 2 2006 15:04:05 -0700 MST"

type entryHeader struct {
	Filepath string    `yaml:"-"`
	Filename string    `yaml:"-"`
	Tags     []string  `yaml:"tags,omitempty"`
	Date     time.Time `yaml:"date,omitempty"`
	Content  string    `fm:"content" yaml:"-"`
}

func (e *entryHeader) MarshalFrontmatter() ([]byte, error) {
	return frontmatter.Marshal(&struct {
		Tags    []string `yaml:"tags,omitempty"`
		Date    string   `yaml:"date"`
		Content string   `fm:"content" yaml:"-"`
	}{
		Tags:    e.Tags,
		Date:    e.Date.Format(JournalTimeformat),
		Content: e.Content,
	})
}

func unmarshalFrontmatter(input []byte) (*entryHeader, error) {
	type rawHeader struct {
		Tags    []string `yaml:"tags,omitempty"`
		Date    string   `yaml:"date,omitempty"`
		Content string   `fm:"content" yaml:"-"`
	}
	raw := new(rawHeader)
	frontmatter.Unmarshal(input, raw)
	date, err := time.Parse(JournalTimeformat, raw.Date)
	if err != nil {
		return nil, err
	}
	return &entryHeader{
		Tags:    raw.Tags,
		Date:    date,
		Content: raw.Content,
	}, nil
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
	head, err := unmarshalFrontmatter(content)
	head.Filepath = filePath
	head.Filename = path.Base(filePath)
	results <- frontmatterResult{
		header: head,
		err:    err,
	}
}
