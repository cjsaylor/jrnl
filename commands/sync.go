package commands

import (
	"context"
	"errors"
)

type SyncCommand struct {
	options Configuration
}

// NewSyncCommand creates a new command runner for sync command
func NewSyncCommand(config Configuration) *SyncCommand {
	syncCommand := SyncCommand{
		options: config,
	}
	return &syncCommand
}

// Run the sync command
func (s *SyncCommand) Run(ctx context.Context, subcommandArgs []string) error {
	params := []string{
		"-C",
		s.options.JournalPath,
		"pull",
	}
	if code := gitCommand(params...); code != 0 {
		return errors.New("failed to sync journal")
	}
	return nil
}
