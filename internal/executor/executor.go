package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"forge/internal/template"
)

// Executor runs commands in an isolated workspace
type Executor struct {
	workDir     string
	interactive bool
}

// New creates a new command executor
func New(workDir string, interactive bool) *Executor {
	return &Executor{
		workDir:     workDir,
		interactive: interactive,
	}
}

// Run executes a command in the workspace
func (e *Executor) Run(cmd template.Command) error {
	if len(cmd.Cmd) == 0 {
		return fmt.Errorf("empty command")
	}
	
	// Create command
	execCmd := exec.Command(cmd.Cmd[0], cmd.Cmd[1:]...)
	execCmd.Dir = e.workDir
	
	if e.interactive {
		// Interactive mode: connect stdin/stdout/stderr to user terminal
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
	} else {
		// Non-interactive mode: close stdin, capture output
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
			
			// Check if this might be an interactive command
			if isLikelyInteractiveError(stderr.String()) {
				return fmt.Errorf("%w\n\nHint: This command may require user input. Try running with --interactive flag", err)
			}
			
			return err
		}
		
		return nil
	}
	
	// Run in interactive mode
	return execCmd.Run()
}

// isLikelyInteractiveError checks if an error message suggests the command wanted interactive input
func isLikelyInteractiveError(stderr string) bool {
	// Common patterns that suggest interactive input was expected
	indicators := []string{
		"stdin",
		"input",
		"prompt",
		"interactive",
		"tty",
		"terminal",
	}
	
	stderrLower := stderr
	for _, indicator := range indicators {
		if len(stderrLower) > 0 && strings.Contains(stderrLower, indicator) {
			return true
		}
	}
	
	return false
}
