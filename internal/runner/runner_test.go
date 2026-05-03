package runner

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunCapturesStdoutStderrAndMultilineOutput(t *testing.T) {
	requireShell(t)

	result, err := Run(context.Background(), "sh", "-c", "printf 'one\n'; printf 'two\n' >&2")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0", result.ExitCode)
	}

	got := plainLines(result)
	if !contains(got, "one") || !contains(got, "two") {
		t.Fatalf("captured lines = %#v, want stdout and stderr lines", got)
	}
}

func TestRunReturnsReviewableNonZeroExit(t *testing.T) {
	requireShell(t)

	result, err := Run(context.Background(), "sh", "-c", "printf 'failed\n'; exit 7")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.ExitCode != 7 {
		t.Fatalf("exit code = %d, want 7", result.ExitCode)
	}
	if got := plainLines(result); !contains(got, "failed") {
		t.Fatalf("captured lines = %#v, want %q", got, "failed")
	}
}

func TestRunMissingCommandReturnsError(t *testing.T) {
	_, err := Run(context.Background(), "rvw-command-that-does-not-exist")
	if err == nil {
		t.Fatal("Run() error = nil, want error")
	}
}

func TestRunHonorsContextTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PTY runner is not supported on Windows yet")
	}
	requireShell(t)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := Run(ctx, "sh", "-c", "sleep 5")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run() error = %v, want context deadline exceeded", err)
	}
}

func TestRunTruncatesLargeOutput(t *testing.T) {
	requireShell(t)

	result, err := RunWithOptions(context.Background(), Options{MaxOutputBytes: 5}, "sh", "-c", "printf '1234567890'")
	if err != nil {
		t.Fatalf("RunWithOptions() error = %v", err)
	}
	if !result.Truncated {
		t.Fatal("truncated = false, want true")
	}

	joined := strings.Join(plainLines(result), "\n")
	if !strings.Contains(joined, "[rvw: output truncated after 5 bytes]") {
		t.Fatalf("captured output = %q, want truncation marker", joined)
	}
}

func plainLines(result Result) []string {
	lines := make([]string, 0, len(result.Transcript.Lines))
	for _, line := range result.Transcript.Lines {
		lines = append(lines, line.Plain)
	}
	return lines
}

func contains(lines []string, want string) bool {
	for _, line := range lines {
		if line == want {
			return true
		}
	}
	return false
}

func requireShell(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell runner tests require sh")
	}
}
