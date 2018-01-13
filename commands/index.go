package commands

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/ericaro/frontmatter"
)

type IndexCommand struct {
	options Configuration
	flags   *flag.FlagSet
}

type entryHeader struct {
	Filename string
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
	head.Filename = path.Base(filePath)
	results <- frontmatterResult{
		header: head,
		err:    nil,
	}
}

func tagMap(journalPath string) (map[string][]string, error) {
	directory := journalPath + "/entries"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	index := make(map[string][]string)
	results := make(chan frontmatterResult, len(files))
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			readFrontmatter(filePath, results)
		}(fmt.Sprintf("%s/%s", directory, file.Name()))
	}
	wg.Wait()
	close(results)
	for result := range results {
		if result.err != nil {
			return nil, result.err
		}
		for _, tag := range result.header.Tags {
			index[tag] = append(index[tag], strings.TrimSuffix(result.header.Filename, ".md"))
		}
	}
	for tag := range index {
		sort.Strings(index[tag])
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
		flags:   flag.NewFlagSet("index", flag.ExitOnError),
	}
	return &indexCommand
}

// Run the index command
func (i *IndexCommand) Run(ctx context.Context, subcommandArgs []string) {
	outputPath := i.flags.String("o", "Index.md", "Output path contained to the $JOURNAL_PATH.")
	if !i.flags.Parsed() {
		i.flags.Parse(subcommandArgs)
	}
	if *outputPath == "." {
		*outputPath = "Index.md"
	}
	index, err := tagMap(i.options.JournalPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
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
	indexPath := fmt.Sprintf("%s/%s", i.options.JournalPath, path.Base(*outputPath))
	err = ioutil.WriteFile(indexPath, []byte(newIndex), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
