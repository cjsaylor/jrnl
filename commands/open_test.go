package commands_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cjsaylor/jrnl/commands"
)

type fakeEditor struct {
	called map[string][]string
}

func (f *fakeEditor) OpenEditor(editor string, args ...string) error {
	f.called["open_editor"] = append([]string{editor}, args...)
	return nil
}

func TestFileCreatedOnStartup(t *testing.T) {
	path, _ := filepath.Abs("../fixtures")
	expectedFilePath := path + "/entries/2018-07-28.md"
	t.Run("FileCreatedOnStartup", func(t *testing.T) {
		config := commands.Configuration{
			JournalPath:   path,
			JournalEditor: "vim",
		}
		editor := fakeEditor{
			called: make(map[string][]string),
		}
		cmd := commands.NewOpenCommand(config, &editor)
		ctx := context.WithValue(context.Background(), commands.CommandContextKey("date"), time.Date(2018, time.July, 28, 0, 0, 0, 0, time.UTC))
		if err := cmd.Run(ctx, []string{}); err != nil {
			t.Fatal(err)
		}
		content, err := ioutil.ReadFile(expectedFilePath)
		if err != nil {
			t.Fatal(err)
		}
		expectedContent := "---\ndate: Sat Jul 28 2018 00:00:00 +0000 UTC\n---\n"
		if string(content) != expectedContent {
			t.Errorf("Expected %v, got %v", expectedContent, string(content))
		}
		if len(editor.called["open_editor"]) != 2 {
			t.Error("Expected OpenEditor to be called properly.")
		}
		if editor.called["open_editor"][0] != "vim" || editor.called["open_editor"][1] != expectedFilePath {
			t.Errorf("Expected editor to be vim, got %v", editor.called["open_editor"][0])
			t.Errorf("Expected file path to be %v, got %v", expectedFilePath, editor.called["open_editor"][1])
		}
	})
	os.Remove(expectedFilePath)
}
