package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/ericaro/frontmatter"
)

type configuration struct {
	JournalPath          string `env:"JOURNAL_PATH"`
	JournalEditor        string `env:"JRNL_EDITOR" envDefault:"vim"`
	JournalEditorOptions string `env:"JRNL_EDITOR_OPTIONS"`
}

var config configuration

type entryHeader struct {
	Tags    []string `yaml:"tags"`
	Content string   `fm:"content" yaml:"-"`
}

var availableCommands = map[string]string{
	"open":     "    Open a journal entry in configured editor.",
	"memorize": "Commit all journal entries.",
	"sync":     "    Syncronize journal entries from source.",
	"index":    "   Write index file based on frontmatter tags.",
	"image":    "  Append an image to the current journal entry.",
}

func init() {
	config = configuration{}
	env.Parse(&config)
	if config.JournalPath == "" {
		config.JournalPath = os.Getenv("HOME") + "/journal.wiki"
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: jrnl [options...] [command]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n\n")
		for key, description := range availableCommands {
			fmt.Fprintf(os.Stderr, "%s %s\n", key, description)
		}
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
}

func getEnv(key string, defaultValue string) string {
	if value, existed := os.LookupEnv(key); existed {
		return value
	}
	return defaultValue
}

func gitCommand(params ...string) int {
	cmd := exec.Command("git", params...)
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
	}
	return 0
}

func checkErr(err error, code int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
}

func main() {
	now := time.Now()

	date := flag.String("date", now.Format("2006-01-02"), "Specify the date of entry.")
	flag.Parse()

	commands := flag.Args()
	var command string
	if len(commands) < 1 {
		command = "open"
	} else {
		command = commands[0]
		commands = commands[1:]
	}
	switch command {
	case "open":
		var options []string
		if editorOptions := config.JournalEditorOptions; editorOptions != "" {
			options = strings.Split(editorOptions, " ")
		}
		options = append(options, fmt.Sprintf("%s/entries/%s.md", config.JournalPath, *date))
		os.MkdirAll(config.JournalPath+"/entries", os.ModePerm)
		cmd := exec.Command(config.JournalEditor, options...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		checkErr(err, 2)
		break
	case "memorize":
		params := []string{
			"-C",
			config.JournalPath,
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
			config.JournalPath,
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
			config.JournalPath,
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
		break
	case "sync":
		params := []string{
			"-C",
			config.JournalPath,
			"pull",
		}
		if code := gitCommand(params...); code != 0 {
			fmt.Fprintln(os.Stderr, "Failed to sync journal.")
			os.Exit(5)
		}
		break
	case "index":
		directory := config.JournalPath + "/entries"
		files, err := ioutil.ReadDir(directory)
		checkErr(err, 6)
		index := make(map[string][]string)
		for _, file := range files {
			content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, file.Name()))
			checkErr(err, 7)
			head := new(entryHeader)
			frontmatter.Unmarshal(content, head)
			for _, tag := range head.Tags {
				index[tag] = append(index[tag], strings.TrimSuffix(file.Name(), ".md"))
			}
		}
		keys := make([]string, 0, len(index))
		for key := range index {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var newIndex string
		for _, tag := range keys {
			newIndex += fmt.Sprintf("\n* *%s* ", tag)
			mappedEntries := make([]string, len(index[tag]))
			mapper := func(entry string) string {
				return fmt.Sprintf("[%s](%s)", entry, entry)
			}
			for i, entry := range index[tag] {
				mappedEntries[i] = mapper(entry)
			}
			newIndex += strings.Join(mappedEntries, ", ")
		}
		err = ioutil.WriteFile(config.JournalPath+"/Index.md", []byte(newIndex), 0644)
		checkErr(err, 7)
		break
	case "image":
		if len(commands) == 0 {
			fmt.Fprintf(os.Stderr, "Must provide file path.")
			os.Exit(8)
		}
		data, err := ioutil.ReadFile(commands[0])
		checkErr(err, 9)
		os.MkdirAll(config.JournalPath+"/bin", os.ModePerm)
		err = ioutil.WriteFile(config.JournalPath+"/bin/"+filepath.Base(commands[0]), data, 0644)
		checkErr(err, 10)
		appendText := `
---

![](bin/%s)`
		journalEntry := config.JournalPath + "/entries/" + *date + ".md"
		f, err := os.OpenFile(journalEntry, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			f, err = os.Create(journalEntry)
		}
		checkErr(err, 11)
		_, err = f.WriteString(fmt.Sprintf(appendText, filepath.Base(commands[0])))
		checkErr(err, 12)
		break
	default:
		fmt.Fprintf(os.Stderr, "Command [%s] not known.", command)
		os.Exit(1)
	}
}
