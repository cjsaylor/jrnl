package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type OpenCommand struct {
	options Configuration
	flags   *flag.FlagSet
}

// NewOpenCommand creates a new command runner for open command
func NewOpenCommand(config Configuration) *OpenCommand {
	openCommand := OpenCommand{
		options: config,
		flags:   flag.NewFlagSet("open", flag.ExitOnError),
	}
	return &openCommand
}

// Run the open command
func (o *OpenCommand) Run(ctx context.Context, subcommandArgs []string) {
	subjectFlag := o.flags.String("s", "", "Set the subject (this will not use a journal date.")
	if !o.flags.Parsed() {
		o.flags.Parse(subcommandArgs)
	}
	var filename string
	if *subjectFlag != "" {
		filename = *subjectFlag
	} else {
		filename = ctx.Value(CommandContextKey("date")).(time.Time).Format("2006-01-02")
	}
	var options []string
	if editorOptions := o.options.JournalEditorOptions; editorOptions != "" {
		options = strings.Split(editorOptions, " ")
	}
	options = append(options, fmt.Sprintf("%s/entries/%s.md", o.options.JournalPath, filename))
	os.MkdirAll(o.options.JournalPath+"/entries", os.ModePerm)
	cmd := exec.Command(o.options.JournalEditor, options...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
