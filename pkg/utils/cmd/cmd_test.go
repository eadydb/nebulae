package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"testing"
)

func helperCommand(s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// adapted from https://npf.io/2015/06/testing-exec-command
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "skaffold":
		var iargs []interface{}
		for _, s := range args {
			iargs = append(iargs, s)
		}
		fmt.Println(iargs...)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(2)
	}

	os.Exit(0)
}

func TestCmd_RunCmdOut(t *testing.T) {
	tests := []struct {
		description string
		cmd         *exec.Cmd
		want        string
		shouldErr   bool
	}{
		{
			description: "skaffold test",
			cmd:         helperCommand("skaffold", "dev"),
			want:        "dev\n",
			shouldErr:   false,
		},
		{
			description: "unknown command test",
			cmd:         helperCommand("foo", "bar"),
			want:        "",
			shouldErr:   true,
		},
	}
	for _, test := range tests {
		got, err := RunCmdOut(context.Background(), test.cmd)
		slog.Info("msg got", slog.String("got", string(got)))
		if test.shouldErr && err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
