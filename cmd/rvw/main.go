package main

import (
	"fmt"
	"os"

	"github.com/V-Isa/rvw/internal/app"
)

func main() {
	if err := app.NewRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
