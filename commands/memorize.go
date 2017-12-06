package commands

import (
	"context"
	"fmt"
	"os"
)

type MemorizeCommand struct {
	options Configuration
}

// NewMemorizeCommand creates a new command runner for memorize command
func NewMemorizeCommand(config Configuration) *MemorizeCommand {
	memorizeCommand := MemorizeCommand{
		options: config,
	}
	return &memorizeCommand
}

// Run the memorize command
func (m *MemorizeCommand) Run(ctx context.Context, subcommandArgs []string) {
	params := []string{
		"-C",
		m.options.JournalPath,
		"add",
		".",
	}

	if code := gitCommand(params...); code != 0 {
		switch code {
		case 128:
			break
		default:
			fmt.Fprintln(os.Stderr, "Failed to stage journal entries.")
			os.Exit(4)
		}
	}

	params = []string{
		"-C",
		m.options.JournalPath,
		"commit",
		"-am",
		"Memorized journal",
	}

	if code := gitCommand(params...); code != 0 {
		switch code {
		case 128:
			break
		default:
			fmt.Fprintln(os.Stderr, "Failed to commit journal.")
			os.Exit(3)
		}
	}

	params = []string{
		"-C",
		m.options.JournalPath,
		"push",
		"origin",
		"master",
	}

	if code := gitCommand(params...); code != 0 {
		switch code {
		case 128:
			break
		default:
			fmt.Fprintln(os.Stderr, "Failed to sync journal.")
			os.Exit(4)
		}
	}
}
