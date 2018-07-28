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
func (t *TagCommand) Run(ctx context.Context, subcommandArgs []string) error {
	var files arrayFlags
	var subjects arrayFlags
	var tags arrayFlags
	var dates arrayFlags
	t.flags.Var(&files, "f", "File path of document to tag")
	t.flags.Var(&subjects, "s", "Subject(s) entries to tag")
	t.flags.Var(&dates, "d", "Specify the date(s) of entry.")
	t.flags.Var(&tags, "t", "Tag or tags to append to specified files, subjects, or dates")
	if !t.flags.Parsed() {
		if err := t.flags.Parse(subcommandArgs); err != nil {
			return err
		}
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
			return err
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
	// @todo Make this async for performance after certain len()
	for result := range results {
		if result.err != nil {
			return result.err
		}
		result.header.Tags = dedupe(append(result.header.Tags, tags...))
		sort.Strings(result.header.Tags)
		output, err := result.header.MarshalFrontmatter()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(result.header.Filepath, output, 0644)
		if err != nil {
			return err
		}
	}
	return nil
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
