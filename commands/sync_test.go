package commands_test

import (
	"context"
	"testing"
	"time"

	"github.com/cjsaylor/jrnl/commands"
)

type fakeGitCommand struct {
	called map[string][]string
}

func (f fakeGitCommand) Pull(path string) error {
	f.called["pull"] = []string{path}
	return nil
}

func TestSyncRun(t *testing.T) {
	runner := fakeGitCommand{
		called: make(map[string][]string),
	}
	cmd := commands.NewSyncCommand(config, runner)
	ctx = context.WithValue(context.Background(), commands.CommandContextKey("date"), time.Date(2018, time.July, 28, 0, 0, 0, 0, time.UTC))
	if err := cmd.Run(ctx, []string{}); err != nil {
		t.Error(err)
	}
	if len(runner.called["pull"]) == 0 {
		t.Error("Expected a git pull")
		return
	}
	expectedPullPath := config.JournalPath
	if runner.called["pull"][0] != expectedPullPath {
		t.Errorf("Expected %v, got %v.", expectedPullPath, runner.called["pull"][0])
	}
}
