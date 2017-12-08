package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/ericaro/frontmatter"
)

type IndexCommand struct {
	options Configuration
}

type entryHeader struct {
	Tags    []string `yaml:"tags"`
	Content string   `fm:"content" yaml:"-"`
}

func tagMap(journalPath string) (map[string][]string, error) {
	directory := journalPath + "/entries"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	index := make(map[string][]string)
	for _, file := range files {
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, file.Name()))
		if err != nil {
			return nil, err
		}
		head := new(entryHeader)
		frontmatter.Unmarshal(content, head)
		for _, tag := range head.Tags {
			index[tag] = append(index[tag], strings.TrimSuffix(file.Name(), ".md"))
		}
	}
	return index, nil
}

func sortedTagKeys(index map[string][]string) []string {
	keys := make([]string, 0, len(index))
	for key := range index {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// NewIndexCommand creates a new command runner for index command
func NewIndexCommand(config Configuration) *IndexCommand {
	indexCommand := IndexCommand{
		options: config,
	}
	return &indexCommand
}

// Run the index command
func (i *IndexCommand) Run(ctx context.Context, subcommandArgs []string) {
	index, err := tagMap(i.options.JournalPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keys := sortedTagKeys(index)
	var newIndex string
	for _, tag := range keys {
		newIndex += fmt.Sprintf("\n* *%s* ", tag)
		mappedEntries := make([]string, len(index[tag]))
		mapper := func(entry string) string {
			return fmt.Sprintf("[%s](%s)", entry, entry)
		}
		for i, entry := range index[tag] {
			mappedEntries[i] = mapper(entry)
		}
		newIndex += strings.Join(mappedEntries, ", ")
	}
	err = ioutil.WriteFile(i.options.JournalPath+"/Index.md", []byte(newIndex), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
