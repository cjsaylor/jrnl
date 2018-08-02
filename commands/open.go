package commands

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type OpenCommand struct {
	options       Configuration
	flags         *flag.FlagSet
	editorSpawner ExternalEditor
}

type ExternalEditor interface {
	OpenEditor(editor string, args ...string) error
}

type ExternalEditorImpl struct{}

func (e *ExternalEditorImpl) OpenEditor(editor string, args ...string) error {
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// NewOpenCommand creates a new command runner for open command
func NewOpenCommand(config Configuration, editorSpawner ExternalEditor) *OpenCommand {
	openCommand := OpenCommand{
		options:       config,
		flags:         flag.NewFlagSet("open", flag.ExitOnError),
		editorSpawner: editorSpawner,
	}
	return &openCommand
}

func generateFrontmatter(ctx context.Context) ([]byte, error) {
	entry := entryHeader{
		Date: ctx.Value(CommandContextKey("date")).(time.Time),
	}
	return entry.MarshalFrontmatter()
}

// Run the open command
func (o *OpenCommand) Run(ctx context.Context, subcommandArgs []string) error {
	subjectFlag := o.flags.String("s", "", "Set the subject (this will not use a journal date.")
	if !o.flags.Parsed() {
		if err := o.flags.Parse(subcommandArgs); err != nil {
			return err
		}
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
	filePath := fmt.Sprintf("%s/entries/%s.md", o.options.JournalPath, filename)
	options = append(options, filePath)
	os.MkdirAll(o.options.JournalPath+"/entries", os.ModePerm)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		content, err := generateFrontmatter(ctx)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
			return err
		}
	}

	return o.editorSpawner.OpenEditor(o.options.JournalEditor, options...)
}
