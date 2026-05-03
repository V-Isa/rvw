package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"

	"github.com/creack/pty"

	"github.com/V-Isa/rvw/internal/transcript"
)

type Result struct {
	Transcript transcript.Transcript
	ExitCode   int
	Truncated  bool
}

type Options struct {
	MaxOutputBytes int64
}

const DefaultMaxOutputBytes int64 = 25 * 1024 * 1024
const truncationMarkerFormat = "[rvw: output truncated after %d bytes]\n"

func Run(ctx context.Context, name string, args ...string) (Result, error) {
	return RunWithOptions(ctx, Options{MaxOutputBytes: DefaultMaxOutputBytes}, name, args...)
}

func RunWithOptions(ctx context.Context, opts Options, name string, args ...string) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if opts.MaxOutputBytes <= 0 {
		opts.MaxOutputBytes = DefaultMaxOutputBytes
	}

	cmd := exec.CommandContext(ctx, name, args...)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return Result{}, fmt.Errorf("start command: %w", err)
	}

	ctxDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = ptmx.Close()
		case <-ctxDone:
		}
	}()

	var buf limitedBuffer
	buf.max = opts.MaxOutputBytes
	_, copyErr := io.Copy(&buf, ptmx)
	close(ctxDone)
	closeErr := ptmx.Close()
	waitErr := cmd.Wait()

	if ctxErr := ctx.Err(); ctxErr != nil {
		return Result{}, fmt.Errorf("command canceled: %w", ctxErr)
	}
	if copyErr != nil && !isTerminalReadDone(copyErr) {
		return Result{}, fmt.Errorf("capture command output: %w", copyErr)
	}
	if closeErr != nil && waitErr == nil {
		return Result{}, fmt.Errorf("close pty: %w", closeErr)
	}

	exitCode := 0
	if waitErr != nil {
		var exitErr *exec.ExitError
		if !errors.As(waitErr, &exitErr) {
			return Result{}, fmt.Errorf("wait for command: %w", waitErr)
		}
		exitCode = exitErr.ExitCode()
	}

	return Result{
		Transcript: transcript.FromBytes(buf.transcriptBytes()),
		ExitCode:   exitCode,
		Truncated:  buf.truncated,
	}, nil
}

func isTerminalReadDone(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EIO)
}

type limitedBuffer struct {
	buf       bytes.Buffer
	max       int64
	truncated bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	remaining := b.max - int64(b.buf.Len())
	if remaining > 0 {
		n := min(int64(len(p)), remaining)
		_, _ = b.buf.Write(p[:n])
		if n < int64(len(p)) {
			b.truncated = true
		}
	} else if len(p) > 0 {
		b.truncated = true
	}
	return len(p), nil
}

func (b *limitedBuffer) transcriptBytes() []byte {
	out := b.buf.Bytes()
	if !b.truncated {
		return out
	}

	var result bytes.Buffer
	result.Write(out)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		result.WriteByte('\n')
	}
	_, _ = fmt.Fprintf(&result, truncationMarkerFormat, b.max)
	return result.Bytes()
}
