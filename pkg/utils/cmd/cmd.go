package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type cmdError struct {
	args   []string
	stdout []byte
	stderr []byte
	cause  error
}

func (e *cmdError) Error() string {
	return fmt.Sprintf("running %s\n - stdout: %q\n - stderr: %q\n - cause: %s", e.args, e.stdout, e.stderr, e.cause)
}

func (e *cmdError) Unwrap() error {
	return e.cause
}

func (e *cmdError) ExitCode() int {
	if exitError, ok := e.cause.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 0
}

// DefaultExecCommand runs commands using exec.Cmd
var DefaultExecCommand Command = newCommander()

func newCommander() *Commander {
	return &Commander{
		store: NewSyncStore[[]byte](),
	}
}

// Command is an interface used to run commands. All packages should use this
// interface instead of calling exec.Cmd directly.
type Command interface {
	RunCmdOut(ctx context.Context, cmd *exec.Cmd) ([]byte, error)
	RunCmd(ctx context.Context, cmd *exec.Cmd) error
	RunCmdOutOnce(ctx context.Context, cmd *exec.Cmd) ([]byte, error)
}

func RunCmdOut(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	return DefaultExecCommand.RunCmdOut(ctx, cmd)
}

func RunCmd(ctx context.Context, cmd *exec.Cmd) error {
	return DefaultExecCommand.RunCmd(ctx, cmd)
}

func RunCmdOutOnce(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	return DefaultExecCommand.RunCmdOutOnce(ctx, cmd)
}

// Commander is the exec.Cmd implementation of the Command interface
type Commander struct {
	store *SyncStore[[]byte]
}

// RunCmdOut runs an exec.Command and returns the stdout and error.
func (*Commander) RunCmdOut(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	slog.Info("Running command", "args", slog.String("args", strings.Join(cmd.Args, ",")))

	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout
	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting command %v: %w", cmd, err)
	}

	if err := cmd.Wait(); err != nil {
		return stdout.Bytes(), &cmdError{
			args:   cmd.Args,
			stdout: stdout.Bytes(),
			stderr: stderr.Bytes(),
			cause:  err,
		}
	}

	if stderr.Len() > 0 {
		slog.Info("Command output", slog.String("stdout", stdout.String()), slog.String("error", stderr.String()))
	} else {
		slog.Info("Command output", slog.String("stdout", stdout.String()))
	}

	return stdout.Bytes(), nil
}

// RunCmd runs an exec.Command.
func (*Commander) RunCmd(ctx context.Context, cmd *exec.Cmd) error {
	slog.Info("Running command", slog.String("args", strings.Join(cmd.Args, ",")))
	return cmd.Run()
}

func (c *Commander) RunCmdOutOnce(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	return c.store.Exec(cmd.String(), func() ([]byte, error) {
		return RunCmdOut(ctx, cmd)
	})
}
