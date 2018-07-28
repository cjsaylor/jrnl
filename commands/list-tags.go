package commands

import "context"
import "fmt"
import "os"

type ListTagsCommand struct {
	options Configuration
}

// NewListTagsCommand creates a new command runner for listing tags.
func NewListTagsCommand(config Configuration) *ListTagsCommand {
	listTagsCommand := ListTagsCommand{
		options: config,
	}
	return &listTagsCommand
}

// Run the list-tags command
func (l *ListTagsCommand) Run(ctd context.Context, subcommandArgs []string) error {
	index, err := tagMap(l.options.JournalPath)
	if err != nil {
		return err
	}
	tags := sortedTagKeys(index)
	for _, tag := range tags {
		fmt.Fprintf(os.Stdout, "%s\n", tag)
	}
	return nil
}
