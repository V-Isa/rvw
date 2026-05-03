package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/V-Isa/rvw/internal/clipboard"
	"github.com/V-Isa/rvw/internal/comments"
	"github.com/V-Isa/rvw/internal/exporter"
	"github.com/V-Isa/rvw/internal/runner"
	"github.com/V-Isa/rvw/internal/tui"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

const (
	appName                 = "rvw"
	runCommandUse           = "run <command> [args...]"
	timeoutFlag             = "--timeout"
	timeoutFlagWithAssign   = timeoutFlag + "="
	argumentSeparator       = "--"
	runStatusMessage        = "rvw: running %s; review UI opens after the command exits\n"
	missingCommandError     = "missing command to run"
	missingTimeoutError     = "missing duration for --timeout"
	parseTimeoutErrorPrefix = "parse --timeout"
)

type runFunc func(context.Context, string, ...string) (runner.Result, error)
type reviewFunc func(tui.Model) error

type commandDeps struct {
	run    runFunc
	review reviewFunc
}

func NewRootCommand() *cobra.Command {
	return newRootCommand(commandDeps{
		run: runner.Run,
		review: func(model tui.Model) error {
			_, err := tea.NewProgram(model).Run()
			return err
		},
	})
}

func newRootCommand(deps commandDeps) *cobra.Command {
	root := &cobra.Command{
		Use:           appName,
		Short:         "Review terminal command output",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newRunCommand(deps))
	return root
}

func newRunCommand(deps commandDeps) *cobra.Command {
	return &cobra.Command{
		Use:                runCommandUse,
		Short:              "Run a command and review its captured output",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := parseRunArgs(args)
			if err != nil {
				return err
			}
			if len(opts.args) == 0 {
				return errors.New(missingCommandError)
			}

			command := DisplayCommand(opts.args)
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), runStatusMessage, command)

			ctx := cmd.Context()
			var cancel context.CancelFunc
			if opts.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, opts.timeout)
				defer cancel()
			}

			result, err := deps.run(ctx, opts.args[0], opts.args[1:]...)
			if err != nil {
				return err
			}

			model := tui.NewModel(tui.Config{
				Command:   command,
				ExitCode:  result.ExitCode,
				Lines:     result.Transcript.Lines,
				Comments:  comments.NewStore(),
				Exporter:  exporter.Markdown{},
				Clipboard: clipboard.New(),
			})

			return deps.review(model)
		},
	}
}

type runOptions struct {
	timeout time.Duration
	args    []string
}

func parseRunArgs(args []string) (runOptions, error) {
	var opts runOptions
	for len(args) > 0 {
		arg := args[0]
		switch {
		case arg == argumentSeparator:
			opts.args = args[1:]
			return opts, nil
		case arg == timeoutFlag:
			if len(args) < 2 {
				return opts, errors.New(missingTimeoutError)
			}
			timeout, err := time.ParseDuration(args[1])
			if err != nil {
				return opts, fmt.Errorf("%s: %w", parseTimeoutErrorPrefix, err)
			}
			opts.timeout = timeout
			args = args[2:]
		case strings.HasPrefix(arg, timeoutFlagWithAssign):
			timeout, err := time.ParseDuration(strings.TrimPrefix(arg, timeoutFlagWithAssign))
			if err != nil {
				return opts, fmt.Errorf("%s: %w", parseTimeoutErrorPrefix, err)
			}
			opts.timeout = timeout
			args = args[1:]
		default:
			opts.args = args
			return opts, nil
		}
	}
	return opts, nil
}
