package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
)

const (
	darwinOS      = "darwin"
	pbcopyCommand = "pbcopy"
)

type Clipboard interface {
	WriteText(text string) error
}

func New() Clipboard {
	if runtime.GOOS == darwinOS {
		return macOS{}
	}
	return unsupported{goos: runtime.GOOS}
}

type macOS struct{}

func (macOS) WriteText(text string) error {
	cmd := exec.Command(pbcopyCommand)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("open pbcopy stdin: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start pbcopy: %w", err)
	}
	if _, err := stdin.Write([]byte(text)); err != nil {
		_ = stdin.Close()
		return fmt.Errorf("write clipboard text: %w", err)
	}
	if err := stdin.Close(); err != nil {
		return fmt.Errorf("close clipboard input: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("copy to clipboard: %w", err)
	}
	return nil
}

type unsupported struct {
	goos string
}

func (u unsupported) WriteText(string) error {
	return fmt.Errorf("clipboard unsupported on %s", u.goos)
}
