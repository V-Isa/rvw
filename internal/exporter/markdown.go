package exporter

import (
	"fmt"
	"strings"

	"github.com/V-Isa/rvw/internal/comments"
	"github.com/V-Isa/rvw/internal/transcript"
)

type Markdown struct{}

type Review struct {
	Command  string
	ExitCode int
	Lines    []transcript.Line
	Comments []comments.Comment
}

func (Markdown) Export(review Review) string {
	var b strings.Builder
	b.WriteString("# rvw review\n\n")
	b.WriteString("Command:\n")
	_, _ = fmt.Fprintf(&b, "`%s`\n\n", review.Command)
	_, _ = fmt.Fprintf(&b, "Exit code: `%d`\n\n", review.ExitCode)
	b.WriteString("## Comments\n")

	if len(review.Comments) == 0 {
		b.WriteString("\nNo comments.\n")
		return b.String()
	}

	lineText := make(map[int]string, len(review.Lines))
	for _, line := range review.Lines {
		lineText[line.Number] = line.Plain
	}

	for _, comment := range review.Comments {
		_, _ = fmt.Fprintf(&b, "\n### L%d\n", comment.Line)
		b.WriteString("Selected line:\n")
		_, _ = fmt.Fprintf(&b, "> %s\n\n", quoteLine(lineText[comment.Line]))
		b.WriteString("Comment:\n")
		_, _ = fmt.Fprintf(&b, "%s\n", comment.Text)
	}

	return b.String()
}

func quoteLine(s string) string {
	if s == "" {
		return "_empty line_"
	}
	return strings.ReplaceAll(s, "\n", " ")
}
