package main

import (
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
)

type ExecuteOptions struct {
	TrimSpace    bool
	IncludeStack bool
	Stdin        *string
	// For example "MY_VAR=some_value"
	EnvironmentVariables []string
}

func Execute(programName string, args ...string) string {
	return ExecuteWithOptions(ExecuteOptions{TrimSpace: true, IncludeStack: true}, programName, args...)
}

func ExecuteWithOptions(options ExecuteOptions, programName string, args ...string) string {
	cmd := exec.Command(programName, args...)
	if options.EnvironmentVariables != nil {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, options.EnvironmentVariables...)
	}
	if options.Stdin != nil {
		cmd.Stdin = strings.NewReader(*options.Stdin)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if options.IncludeStack {
			debug.PrintStack()
		}
		log.Fatal("Failed executing ", programName, args, ": ", string(out), err)
	}
	if options.TrimSpace {
		return strings.TrimSpace(string(out))
	} else {
		return string(out)
	}
}

func ExecuteFailable(programName string, args ...string) (string, error) {
	cmd := exec.Command(programName, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
