package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Check defines a single check.
type Check struct {
	Name        string
	Description string
	Run         func() error
}

// AllChecks contains a list of all checks with descriptions.
var AllChecks = []Check{
	{
		Name:        "check-branch-master",
		Description: "test if the current branch is master",
		Run:         CheckBranchMaster,
	},
}

// CheckResult bundles a check with its result after running.
type CheckResult struct {
	Check
	Result error
}

// ErrInvalidCheckName is returned when an unknown name for a check is passed.
type ErrInvalidCheckName struct {
	Name string
}

func (err ErrInvalidCheckName) Error() string {
	return fmt.Sprintf("invalid check name %q", err.Name)
}

// FilterChecks returns a list of checks without the ones listed in reject. For
// invalid names, an error is returned.
func FilterChecks(list []Check, reject []string) (result []Check, err error) {
	all := make(map[string]struct{})
	for _, check := range list {
		all[check.Name] = struct{}{}
	}

	disabled := make(map[string]struct{})

	for _, name := range reject {
		if _, ok := all[name]; !ok {
			return nil, ErrInvalidCheckName{Name: name}
		}

		disabled[name] = struct{}{}
	}

	for _, check := range list {
		if _, ok := disabled[check.Name]; ok {
			continue
		}

		result = append(result, check)
	}

	return result, nil
}

// RunChecks runs all checks.
func RunChecks(checks []Check) (result []CheckResult) {
	for _, check := range checks {
		result = append(result, CheckResult{
			Check:  check,
			Result: check.Run(),
		})
	}

	return result
}

type BranchError struct {
	CurrentBranchName string
}

func (err BranchError) Error() string {
	return fmt.Sprintf("current branch is %q instead of master", err.CurrentBranchName)
}

func CheckBranchMaster() error {
	name, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return fmt.Errorf("unable to find current branch: %w", err)
	}

	branch := strings.TrimRight(string(name), "\n")
	if branch != "master" {
		return BranchError{CurrentBranchName: branch}
	}

	return nil
}
