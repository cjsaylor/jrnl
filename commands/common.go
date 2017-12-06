package commands

import (
	"context"
	"os/exec"
	"syscall"
)

// CommandRunner is an interface for runnable commands
type CommandRunner interface {
	Run(context.Context, []string)
}

// CommandContextKey is a context key specific to commands package
type CommandContextKey string

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
