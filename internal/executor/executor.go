package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"forge/internal/template"
)

// Executor runs commands in a workspace
type Executor struct {
	workDir  string
	testMode bool
}

// New creates a new command executor
func New(workDir string, interactive bool, testMode bool) *Executor {
	// Note: interactive parameter is kept for backward compatibility but ignored
	// forge init: always uses real TTY
	// forge test: always non-interactive
	return &Executor{
		workDir:  workDir,
		testMode: testMode,
	}
}

// Run executes a command in the workspace
func (e *Executor) Run(cmd template.Command) error {
	if len(cmd.Cmd) == 0 {
		return fmt.Errorf("empty command")
	}

	// Determine which command to run
	cmdToRun := cmd.Cmd

	// During test mode, handle interactive commands
	if e.testMode && cmd.Interactive {
		if len(cmd.TestCmd) > 0 {
			// Use test command
			fmt.Printf("[forge test] Using test command for interactive step: %s\n", strings.Join(cmd.TestCmd, " "))
			cmdToRun = cmd.TestCmd
		} else {
			// Skip with warning
			fmt.Printf("[forge test] Skipping interactive command: %s\n", strings.Join(cmd.Cmd, " "))
			return nil
		}
	}

	// Create command
	execCmd := exec.Command(cmdToRun[0], cmdToRun[1:]...)
	execCmd.Dir = e.workDir

	// For forge init: always use real TTY (inherit terminal I/O)
	// For forge test: capture output (never interactive)
	if !e.testMode {
		// forge init mode: connect stdin/stdout/stderr to user terminal
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		return execCmd.Run()
	}

	// forge test mode: non-interactive, capture output
	execCmd.Stdin = nil

	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	// Run command
	if err := execCmd.Run(); err != nil {
		// On error, show captured output
		if stdout.Len() > 0 {
			fmt.Fprintf(os.Stderr, "\nStdout:\n%s\n", stdout.String())
		}
		if stderr.Len() > 0 {
			fmt.Fprintf(os.Stderr, "\nStderr:\n%s\n", stderr.String())
		}

		return err
	}

	return nil
}
