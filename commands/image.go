package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var appendTemplate = `

---

![](bin/%s)`

type ImageCommand struct {
	options Configuration
	flags   *flag.FlagSet
}

// NewImageCommand creates a new command runner for image command
func NewImageCommand(config Configuration) *ImageCommand {
	imageCommand := ImageCommand{
		options: config,
		flags:   flag.NewFlagSet("image", flag.ExitOnError),
	}
	return &imageCommand
}

// Run the image command
func (i *ImageCommand) Run(ctx context.Context, subcommandArgs []string) error {
	subjectFlag := i.flags.String("s", "", "Set the subject (this will not use a journal date.")
	if !i.flags.Parsed() {
		if err := i.flags.Parse(subcommandArgs); err != nil {
			return err
		}
	}
	var filebase string
	if *subjectFlag != "" {
		filebase = *subjectFlag
	} else {
		filebase = ctx.Value(CommandContextKey("date")).(time.Time).Format("2006-01-02")
	}
	commandArgs := i.flags.Args()
	if len(commandArgs) == 0 {
		return errors.New("must provide file path")
	}
	data, err := ioutil.ReadFile(commandArgs[0])
	if err != nil {
		return err
	}
	os.MkdirAll(i.options.JournalPath+"/bin", os.ModePerm)
	err = ioutil.WriteFile(i.options.JournalPath+"/bin/"+filepath.Base(commandArgs[0]), data, 0644)
	if err != nil {
		return err
	}
	journalEntry := i.options.JournalPath + "/entries/" + filebase + ".md"
	f, err := os.OpenFile(journalEntry, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		f, err = os.Create(journalEntry)
	}
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf(appendTemplate, filepath.Base(commandArgs[0])))
	return err
}
