package commands

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ericaro/frontmatter"
)

type TagCommand struct {
	options Configuration
	flags   *flag.FlagSet
}

// NewTagCommand creates a new command runner for tagging entries
func NewTagCommand(config Configuration) *TagCommand {
	tagCommand := TagCommand{
		options: config,
		flags:   flag.NewFlagSet("tag", flag.ExitOnError),
	}
	return &tagCommand
}

// Run the tag command
func (t *TagCommand) Run(ctx context.Context, subcommandArgs []string) {
	var files arrayFlags
	var subjects arrayFlags
	var tags arrayFlags
	var dates arrayFlags
	t.flags.Var(&files, "f", "File path of document to tag")
	t.flags.Var(&subjects, "s", "Subject(s) entries to tag")
	t.flags.Var(&dates, "d", "Specify the date(s) of entry.")
	t.flags.Var(&tags, "t", "Tag or tags to append to specified files, subjects, or dates")
	if !t.flags.Parsed() {
		t.flags.Parse(subcommandArgs)
	}
	var fileEntries []string
	for _, file := range files {
		fileEntries = append(fileEntries, file)
	}
	for _, subject := range subjects {
		fileEntries = append(fileEntries, fmt.Sprintf("%s/entries/%s.md", t.options.JournalPath, subject))
	}
	for _, date := range dates {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fileEntries = append(fileEntries, fmt.Sprintf("%s/entries/%s.md", t.options.JournalPath, parsedDate.String()))
	}
	if len(fileEntries) == 0 {
		toCreate := fmt.Sprintf("%s/entries/%s.md", t.options.JournalPath, ctx.Value(CommandContextKey("date")).(time.Time).Format("2006-01-02"))
		os.OpenFile(toCreate, os.O_RDONLY|os.O_CREATE, 0644)
		fileEntries = append(fileEntries, toCreate)
	}
	var wg sync.WaitGroup
	results := make(chan frontmatterResult, len(fileEntries))
	for _, file := range fileEntries {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			readFrontmatter(filePath, results)
		}(file)
	}
	wg.Wait()
	close(results)
	for result := range results {
		if result.err != nil {
			fmt.Println(result.err)
			os.Exit(2)
		}
		result.header.Tags = dedupe(append(result.header.Tags, tags...))
		sort.Strings(result.header.Tags)
		output, err := frontmatter.Marshal(result.header)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		err = ioutil.WriteFile(result.header.Filepath, output, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	}
}

func dedupe(subject []string) []string {
	encountered := make(map[string]struct{})
	results := []string{}
	for _, str := range subject {
		if _, ok := encountered[str]; !ok {
			results = append(results, str)
			encountered[str] = struct{}{}
		}
	}
	return results
}
