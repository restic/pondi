package main

import (
	"fmt"
	"os/exec"
)

// Hook is a command which is run before releasing.
type Hook struct {
	Name        string
	Description string
	Command     []string
}

var AllHooks = []Hook{
	{
		Name:        "run-go-mod-download",
		Description: "run 'go mod download' to make sure all Go modules are accessible",
		Command:     []string{"go", "mod", "download"},
	},
	{
		Name:        "run-go-generate",
		Description: "run 'go generate ./...' to make sure all generated code is up to date",
		Command:     []string{"go", "generate", "./..."},
	},
}

// RunHooks run all hooks.
func RunHooks() error {
	for _, hook := range AllHooks {
		fmt.Printf("run %v\n", hook.Name)
		cmd := exec.Command(hook.Command[0], hook.Command[1:]...)
		err := cmd.Run()

		if err != nil {
			return fmt.Errorf("hook %v failed: %w", hook.Name, err)
		}
	}

	return nil
}
