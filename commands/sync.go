package commands

import (
	"context"
	"errors"
)

type SyncCommand struct {
	options Configuration
	runner  GitCommandRunner
}

type GitCommandRunner interface {
	Pull(path string) error
}

type GitCommandRunnerImpl struct{}

func (g GitCommandRunnerImpl) Pull(path string) error {
	params := []string{
		"-C",
		path,
		"pull",
	}
	if code := gitCommand(params...); code != 0 {
		return errors.New("failed to sync journal")
	}
	return nil
}

// NewSyncCommand creates a new command runner for sync command
func NewSyncCommand(config Configuration, runner GitCommandRunner) *SyncCommand {
	syncCommand := SyncCommand{
		options: config,
		runner:  runner,
	}
	return &syncCommand
}

// Run the sync command
func (s *SyncCommand) Run(ctx context.Context, subcommandArgs []string) error {
	return s.runner.Pull(s.options.JournalPath)
}
