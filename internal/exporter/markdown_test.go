package exporter

import (
	"strings"
	"testing"

	"github.com/V-Isa/rvw/internal/comments"
	"github.com/V-Isa/rvw/internal/transcript"
)

func TestMarkdownExport(t *testing.T) {
	report := Markdown{}.Export(Review{
		Command:  "echo hello",
		ExitCode: 1,
		Lines: []transcript.Line{
			{Number: 1, Plain: "hello"},
		},
		Comments: []comments.Comment{
			{Line: 1, Text: "review this"},
		},
	})

	for _, want := range []string{
		"# rvw review",
		"`echo hello`",
		"Exit code: `1`",
		"### L1",
		"> hello",
		"review this",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report contains %q = false, want true\nreport:\n%s", want, report)
		}
	}
}

func TestMarkdownExportNoComments(t *testing.T) {
	report := Markdown{}.Export(Review{Command: "true"})
	if !strings.Contains(report, "No comments.") {
		t.Fatalf("report contains empty state = false, want true\nreport:\n%s", report)
	}
}
