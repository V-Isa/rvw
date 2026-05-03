package app

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/V-Isa/rvw/internal/runner"
	"github.com/V-Isa/rvw/internal/tui"

	"github.com/spf13/cobra"
)

func TestRunCommandRequiresCommand(t *testing.T) {
	cmd := newTestRootCommand(t, nil)
	cmd.SetArgs([]string{"run"})

	err := cmd.Execute()
	wantErr := "requires at least 1 arg(s), only received 0"
	if err == nil || err.Error() != wantErr {
		t.Fatalf("Execute() error = %v, want %q", err, wantErr)
	}
}

func TestRunCommandStripsDoubleDash(t *testing.T) {
	var gotName string
	var gotArgs []string
	cmd := newTestRootCommand(t, func(_ context.Context, name string, args ...string) (runner.Result, error) {
		gotName = name
		gotArgs = args
		return runner.Result{}, nil
	})
	cmd.SetArgs([]string{"run", "--", "echo", "hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	wantName := "echo"
	wantArgs := []string{"hello"}
	if gotName != wantName || !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("runner name = %q, args = %#v; want name = %q, args = %#v", gotName, gotArgs, wantName, wantArgs)
	}
}

func TestRunCommandPassesChildFlagsThrough(t *testing.T) {
	var gotArgs []string
	cmd := newTestRootCommand(t, func(_ context.Context, _ string, args ...string) (runner.Result, error) {
		gotArgs = args
		return runner.Result{}, nil
	})
	cmd.SetArgs([]string{"run", "echo", "--child-flag"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	wantArgs := []string{"--child-flag"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("runner args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestRunCommandParsesTimeout(t *testing.T) {
	var hasDeadline bool
	cmd := newTestRootCommand(t, func(ctx context.Context, _ string, _ ...string) (runner.Result, error) {
		_, hasDeadline = ctx.Deadline()
		return runner.Result{}, nil
	})
	cmd.SetArgs([]string{"run", "--timeout=1s", "--", "echo", "hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !hasDeadline {
		t.Fatal("runner context deadline = false, want true")
	}
}

func TestDisplayCommandQuotesShellSensitiveArgs(t *testing.T) {
	got := DisplayCommand([]string{"echo", "hello world", "it's", ""})
	want := "echo 'hello world' 'it'\\''s' ''"
	if got != want {
		t.Fatalf("DisplayCommand() = %q, want %q", got, want)
	}
}

func newTestRootCommand(t *testing.T, run runFunc) *cobra.Command {
	t.Helper()
	if run == nil {
		run = func(context.Context, string, ...string) (runner.Result, error) {
			return runner.Result{}, nil
		}
	}
	cmd := newRootCommand(commandDeps{
		run: run,
		review: func(tui.Model) error {
			return nil
		},
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetContext(context.Background())
	return cmd
}
