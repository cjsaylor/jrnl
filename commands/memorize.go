package commands

import (
	"context"
	"errors"
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
func (m *MemorizeCommand) Run(ctx context.Context, subcommandArgs []string) error {
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
			return errors.New("failed to stage journal entries")
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
			return errors.New("failed to commit journal")
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
			return errors.New("failed to sync journal")
		}
	}
	return nil
}
