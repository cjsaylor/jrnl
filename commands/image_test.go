package commands_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cjsaylor/jrnl/commands"
)

func TestAppendImage(t *testing.T) {
	path, _ := filepath.Abs("../fixtures")
	expectedImagePath := fmt.Sprintf("%v/bin/test-pixel.png", path)
	expectedEntryPath := fmt.Sprintf("%v/entries/2018-07-01.md", path)
	t.Run("appendImage", func(t *testing.T) {
		config := commands.Configuration{
			JournalPath: path,
		}
		cmd := commands.NewImageCommand(config)
		ctx := context.WithValue(context.Background(), commands.CommandContextKey("date"), time.Date(2018, time.July, 1, 0, 0, 0, 0, time.UTC))
		if err := cmd.Run(ctx, []string{fmt.Sprintf("%v/%v", path, "test-pixel.png")}); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(expectedImagePath); os.IsNotExist(err) {
			t.Errorf("Expected image to be copied to %v", expectedImagePath)
		}
		if _, err := os.Stat(expectedEntryPath); os.IsNotExist(err) {
			t.Errorf("Expected entry to be created at %v", expectedImagePath)
		}
		entryContent, err := ioutil.ReadFile(expectedEntryPath)
		if err != nil {
			t.Fatal(err)
		}
		expectedContent := "\n\n---\n\n![](bin/test-pixel.png)\n"
		if string(entryContent) != expectedContent {
			t.Errorf("Expected %v, got %v", expectedContent, string(entryContent))
		}
	})
	os.RemoveAll(fmt.Sprintf("%v/bin/", path))
	os.Remove(fmt.Sprintf("%v/entries/2018-07-01.md", path))
}
