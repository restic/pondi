package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func addCommandCheck(root *cobra.Command, gopts *GlobalOptions, cfg *Config) {
	var cmd = &cobra.Command{
		Use:           "check",
		Short:         "run checks and print the result",
		Long:          "run all checks and print the current result for each",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(*gopts, *cfg, args)
		},
	}

	root.AddCommand(cmd)
}

var ErrCheckFailed = errors.New("check failed")

func runCheck(gopts GlobalOptions, _ Config, _ []string) error {
	checks, err := FilterChecks(AllChecks, gopts.DisableChecks)
	if err != nil {
		return err
	}

	failedChecks := 0

	for _, result := range RunChecks(checks) {
		text := ""
		status := "✓"

		if result.Result != nil {
			text = result.Result.Error()
			status = "✗"
			failedChecks++
		}

		fmt.Printf("%s  %v\t%v\t%s\n", status, result.Check.Name, result.Check.Description, text)
	}

	if failedChecks > 0 {
		return ErrCheckFailed
	}

	return nil
}