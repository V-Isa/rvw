package tui

import (
	"strings"
	"testing"

	"github.com/V-Isa/rvw/internal/comments"
	"github.com/V-Isa/rvw/internal/exporter"
	"github.com/V-Isa/rvw/internal/transcript"

	tea "charm.land/bubbletea/v2"
)

func TestModelNavigationStaysInBounds(t *testing.T) {
	model := testModel()

	model = updateModel(t, model, key("j"))
	model = updateModel(t, model, key("j"))
	if model.cursor != 1 {
		t.Fatalf("cursor = %d, want 1", model.cursor)
	}

	model = updateModel(t, model, key("k"))
	model = updateModel(t, model, key("k"))
	if model.cursor != 0 {
		t.Fatalf("cursor = %d, want 0", model.cursor)
	}
}

func TestModelAddsEditsAndClearsComment(t *testing.T) {
	model := testModel()

	model = updateModel(t, model, key("c"))
	model = updateModel(t, model, key("first"))
	model = updateModel(t, model, specialKey(tea.KeyEnter))
	if got, ok := model.comments.Get(1); !ok || got != "first" {
		t.Fatalf("comment = %q, ok = %v; want comment = %q, ok = true", got, ok, "first")
	}

	model = updateModel(t, model, key("c"))
	model.input = "second"
	model = updateModel(t, model, specialKey(tea.KeyEnter))
	if got, ok := model.comments.Get(1); !ok || got != "second" {
		t.Fatalf("comment = %q, ok = %v; want comment = %q, ok = true", got, ok, "second")
	}

	model = updateModel(t, model, key("c"))
	model.input = ""
	model = updateModel(t, model, specialKey(tea.KeyEnter))
	if _, ok := model.comments.Get(1); ok {
		t.Fatal("comment exists = true, want false")
	}
}

func TestModelExportCopiesMarkdown(t *testing.T) {
	clip := &fakeClipboard{}
	model := testModel()
	model.clipboard = clip
	model.comments.Set(1, "review this")

	model = updateModel(t, model, key("e"))

	if !strings.Contains(clip.text, "# rvw review") || !strings.Contains(clip.text, "review this") {
		t.Fatalf("clipboard text = %q, want Markdown report with comment", clip.text)
	}
	if model.status != "Copied Markdown review to clipboard" {
		t.Fatalf("status = %q, want %q", model.status, "Copied Markdown review to clipboard")
	}
	if got := model.footerText(); !strings.Contains(got, defaultFooterHelp) {
		t.Fatalf("footer = %q, want help text", got)
	}
}

func TestModelExportCopiesMarkdownWithYAlias(t *testing.T) {
	clip := &fakeClipboard{}
	model := testModel()
	model.clipboard = clip
	model.comments.Set(1, "review this")

	model = updateModel(t, model, key("y"))

	if !strings.Contains(clip.text, "review this") {
		t.Fatalf("clipboard text = %q, want Markdown report with comment", clip.text)
	}
	if model.status != "Copied Markdown review to clipboard" {
		t.Fatalf("status = %q, want %q", model.status, "Copied Markdown review to clipboard")
	}
}

func TestRenderCommentsUsesLineAnchorOnly(t *testing.T) {
	model := testModel()
	model.comments.Set(2, "review this")

	got := model.renderComments(40, 3)

	if !strings.Contains(got, "L2: review this") {
		t.Fatalf("comments view = %q, want line-anchored comment", got)
	}
	if strings.Contains(got, "#1") {
		t.Fatalf("comments view = %q, want no sidebar ordinal", got)
	}
}

func testModel() Model {
	return NewModel(Config{
		Command:  "echo hello",
		ExitCode: 0,
		Lines: []transcript.Line{
			{Number: 1, Plain: "hello"},
			{Number: 2, Plain: "world"},
		},
		Comments:  comments.NewStore(),
		Exporter:  exporter.Markdown{},
		Clipboard: &fakeClipboard{},
	})
}

func updateModel(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()
	updated, _ := model.Update(msg)
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next
}

func key(text string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Text: text, Code: []rune(text)[0]})
}

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

type fakeClipboard struct {
	text string
}

func (f *fakeClipboard) WriteText(text string) error {
	f.text = text
	return nil
}
