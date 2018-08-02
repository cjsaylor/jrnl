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

func TestFindTag(t *testing.T) {
	path, _ := filepath.Abs("../fixtures")
	config := commands.Configuration{
		JournalPath: path,
	}
	r, w, _ := os.Pipe()
	cmd := commands.NewFindCommand(config, w)
	ctx := context.WithValue(context.Background(), commands.CommandContextKey("date"), time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC))
	cmd.Run(ctx, []string{"-tag", "foo"})
	w.Close()
	output, _ := ioutil.ReadAll(r)
	expectedOutput := fmt.Sprintf("%v/entries/2018-08-01.md\n", path)
	if expectedOutput != string(output) {
		t.Errorf("Expected %v, got %v", expectedOutput, string(output))
	}
	fmt.Println(string(output))
}
