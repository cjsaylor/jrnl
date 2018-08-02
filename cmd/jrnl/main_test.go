package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var location *time.Location

func init() {
	location, _ = time.LoadLocation("UTC")
	now = time.Date(2018, time.August, 1, 13, 24, 15, 0, location)
}

func TestParseDate(t *testing.T) {
	dateTests := []struct {
		input        string
		output       time.Time
		expectsError bool
	}{
		{"2018-02-28", time.Date(2018, time.February, 28, 13, 24, 15, 0, location), false},
		{"2018-12-31", time.Date(2018, time.December, 31, 13, 24, 15, 0, location), false},
		{"2018-02-31", now, true},
		{"random text", now, true},
	}
	for _, testInput := range dateTests {
		t.Run(testInput.input, func(t *testing.T) {
			output, err := ParseDate(testInput.input)
			if err != nil && !testInput.expectsError {
				t.Fatal(err)
			}
			if err == nil && testInput.expectsError {
				t.Error("Expected input to produce an error")
			} else if err != nil && testInput.expectsError {
				return
			}
			if output != testInput.output {
				t.Errorf("expected %v, got %v", testInput.output, output)
			}
		})
	}
}

func TestLongestStringLength(t *testing.T) {
	inputs := []struct {
		input    []string
		expected int
	}{
		{[]string{"a", "aa", "abc", "bb"}, 3},
		{[]string{"", ""}, 0},
		{[]string{"12345678901234567890", "12345678901234567890"}, 20},
	}
	for _, input := range inputs {
		t.Run(fmt.Sprint(input.input), func(t *testing.T) {
			output := LongestStringLength(input.input)
			if output != input.expected {
				t.Errorf("expected %v, got %v", input.expected, output)
			}
		})
	}
}

func getType(input interface{}) string {
	if t := reflect.TypeOf(input); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func TestFromCommandName(t *testing.T) {
	inputs := []struct {
		input        string
		expected     string
		expectsError bool
	}{
		{"open", "*OpenCommand", false},
		{"memorize", "*MemorizeCommand", false},
		{"sync", "*SyncCommand", false},
		{"index", "*IndexCommand", false},
		{"image", "*ImageCommand", false},
		{"list-tags", "*ListTagsCommand", false},
		{"find", "*FindCommand", false},
		{"tag", "*TagCommand", false},
		{"Unknown", "", true},
	}

	for _, input := range inputs {
		t.Run(input.input, func(t *testing.T) {
			output, err := FromCommandName(input.input)
			if err != nil && !input.expectsError {
				t.Fatal(err)
			} else if err == nil && input.expectsError {
				t.Error("Expected input to produce an error")
			} else if err != nil && input.expectsError {
				return
			}
			outputName := getType(output)
			if outputName != input.expected {
				t.Errorf("expected %v, got %v", input.expected, outputName)
			}
		})
	}
}
