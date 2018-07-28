package commands_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cjsaylor/jrnl/commands"
)

var config commands.Configuration
var called map[string][]string
var ctx context.Context

func TestMain(m *testing.M) {
	config = commands.Configuration{
		JournalPath:   "/some/path",
		JournalEditor: "vim",
	}
	called = make(map[string][]string)
	ctx = context.WithValue(context.Background(), commands.CommandContextKey("date"), time.Date(2018, time.July, 28, 0, 0, 0, 0, time.UTC))
	os.Exit(m.Run())
}

type fakeProducer struct{}

func (f *fakeProducer) EnsureDirectory(path string) {
	called["ensure"] = []string{path}
}
func (f *fakeProducer) InitializeEntry(path string, content []byte) error {
	called["initialize"] = []string{path, string(content[:])}
	return nil
}

type fakeEditor struct{}

func (f *fakeEditor) OpenEditor(editor string, args ...string) error {
	called["open_editor"] = append([]string{editor}, args...)
	return nil
}

func TestFileCreatedOnStartup(t *testing.T) {
	producer := fakeProducer{}
	editor := fakeEditor{}
	cmd := commands.NewOpenCommand(config, &producer, &editor)
	if err := cmd.Run(ctx, []string{}); err != nil {
		t.Error(err)
	}
	// Check that the initialized file is called with correct content
	if len(called["ensure"]) == 0 {
		t.Error("Expected EnsureDirectory to be called.")
	}
	if called["ensure"][0] != "/some/path/entries" {
		t.Errorf("Expected directory path to be /some/path/entries, got %v", called["ensure"][0])
	}
	expectedFilePath := "/some/path/entries/2018-07-28.md"
	expectedContent := "---\ndate: Sat Jul 28 2018 00:00:00 +0000 UTC\n---\n"
	if len(called["initialize"]) == 0 {
		t.Error("Expected InitializeEntry to be called.")
		return
	}
	if called["initialize"][0] != expectedFilePath || called["initialize"][1] != expectedContent {
		t.Errorf("Expected file path to be %v, got %v", expectedFilePath, called["initialize"][0])
		t.Errorf("Expected content to be %v, got %v", expectedContent, called["initialize"][1])
	}
	if len(called["open_editor"]) == 0 {
		t.Error("Expcted OpenEditor to be called.")
	}
	if called["open_editor"][0] != "vim" || called["open_editor"][1] != expectedFilePath {
		t.Errorf("Expected editor to be vim, got %v", called["open_editor"][0])
		t.Errorf("Expected file path to be %v, got %v", expectedFilePath, called["open_editor"][1])
	}
}
