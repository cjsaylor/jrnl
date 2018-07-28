package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/cjsaylor/jrnl/commands"
)

var config commands.Configuration

var availableCommands = map[string]string{
	"open":      "Open a journal entry in configured editor.",
	"memorize":  "Commit all journal entries.",
	"sync":      "Syncronize journal entries from source.",
	"index":     "Write index file based on frontmatter tags.",
	"image":     "Append an image to the current journal entry.",
	"list-tags": "List all tags used in journal entries.",
	"find":      "Find journal entries.",
	"tag":       "Append a tag or tags to journal entries.",
}

var version = "dev"

func longestStringLength(strings []string) int {
	length := 0
	for _, v := range strings {
		if len(v) > length {
			length = len(v)
		}
	}
	return length
}

func init() {
	config = commands.Configuration{}
	env.Parse(&config)
	if config.JournalPath == "" {
		config.JournalPath = os.Getenv("HOME") + "/journal.wiki"
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: jrnl [options...] [command]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n\n")
		var keys []string
		for key := range availableCommands {
			keys = append(keys, key)
		}
		maxLength := longestStringLength(keys)
		sort.Strings(keys)
		for _, command := range keys {
			keyPadding := strings.Repeat(" ", maxLength+1)
			fmt.Fprintf(os.Stderr, "%s %s\n", command+keyPadding[len(command):], availableCommands[command])
		}
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
}

func fromCommandName(name string) (commands.CommandRunner, error) {
	switch name {
	case "open":
		return commands.NewOpenCommand(
			config,
			&commands.FileProducer{},
			&commands.ExternalEditorImpl{}), nil
	case "memorize":
		return commands.NewMemorizeCommand(config), nil
	case "sync":
		return commands.NewSyncCommand(config), nil
	case "index":
		return commands.NewIndexCommand(config), nil
	case "image":
		return commands.NewImageCommand(config), nil
	case "list-tags":
		return commands.NewListTagsCommand(config), nil
	case "find":
		return commands.NewFindCommand(config), nil
	case "tag":
		return commands.NewTagCommand(config), nil
	default:
		return nil, errors.New("Command not found")
	}
}

func main() {
	now := time.Now()
	dateInput := flag.String("date", now.Format("2006-01-02"), "Specify the date of entry.")
	versionRequested := flag.Bool("version", false, "Prints the current version.")
	flag.Parse()

	if *versionRequested {
		fmt.Println(version)
		os.Exit(0)
	}

	parsedDate, err := time.Parse("2006-01-02", *dateInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse date: %v. Must be in form of YYYY-mm-dd", *dateInput)
		os.Exit(1)
	}

	ctx := context.WithValue(context.Background(), commands.CommandContextKey("date"), parsedDate)

	commandArgs := flag.Args()
	var command string
	if len(commandArgs) < 1 {
		command = "open"
	} else {
		command = commandArgs[0]
		commandArgs = commandArgs[1:]
	}
	cmd, err := fromCommandName(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd.Run(ctx, commandArgs)
}