package tui

import (
	"fmt"
	"strings"

	"github.com/V-Isa/rvw/internal/clipboard"
	"github.com/V-Isa/rvw/internal/comments"
	"github.com/V-Isa/rvw/internal/exporter"
	"github.com/V-Isa/rvw/internal/transcript"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	headerHeight         = 1
	footerHeight         = 1
	minBodyHeight        = 1
	outputPanePercent    = 70
	percentScale         = 100
	minOutputPaneWidth   = 20
	minCommentsPaneWidth = 10
	paneGapWidth         = 1
	lineNumberWidth      = 4

	emptyCommentMarker = " "
	lineCommentMarker  = "*"

	defaultFooterHelp = "j/k or up/down: move  c: comment  e/y: export/copy Markdown  q: quit"
)

type Config struct {
	Command   string
	ExitCode  int
	Lines     []transcript.Line
	Comments  *comments.Store
	Exporter  exporter.Markdown
	Clipboard clipboard.Clipboard
}

type Model struct {
	command   string
	exitCode  int
	lines     []transcript.Line
	comments  *comments.Store
	exporter  exporter.Markdown
	clipboard clipboard.Clipboard

	width  int
	height int
	cursor int
	offset int

	editing bool
	input   string
	status  string
}

func NewModel(cfg Config) Model {
	return Model{
		command:   cfg.Command,
		exitCode:  cfg.ExitCode,
		lines:     cfg.Lines,
		comments:  cfg.Comments,
		exporter:  cfg.Exporter,
		clipboard: cfg.Clipboard,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.keepCursorVisible()
		return m, nil
	case tea.KeyMsg:
		if m.editing {
			return m.updateInput(msg)
		}
		return m.updateNormal(msg)
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		if m.cursor < len(m.lines)-1 {
			m.cursor++
			m.status = ""
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
			m.status = ""
		}
	case "c":
		m.editing = true
		m.status = ""
		if len(m.lines) > 0 {
			m.input, _ = m.comments.Get(m.lines[m.cursor].Number)
		}
	case "e", "y":
		m = m.exportReview()
	}
	m = m.keepCursorVisible()
	return m, nil
}

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if len(m.lines) > 0 {
			line := m.lines[m.cursor].Number
			m.comments.Set(line, m.input)
			if strings.TrimSpace(m.input) == "" {
				m.status = fmt.Sprintf("Removed comment on L%d", line)
			} else {
				m.status = fmt.Sprintf("Saved comment on L%d", line)
			}
		}
		m.input = ""
		m.editing = false
	case "esc":
		m.input = ""
		m.editing = false
		m.status = "Comment edit canceled"
	case "backspace", "ctrl+h":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	case "space":
		m.input += " "
	default:
		if text := msg.Key().Text; text != "" {
			m.input += text
		}
	}
	return m, nil
}

func (m Model) exportReview() Model {
	report := m.exporter.Export(exporter.Review{
		Command:  m.command,
		ExitCode: m.exitCode,
		Lines:    m.lines,
		Comments: m.comments.All(),
	})
	if err := m.clipboard.WriteText(report); err != nil {
		m.status = "Clipboard error: " + err.Error()
		return m
	}
	m.status = "Copied Markdown review to clipboard"
	return m
}

func (m Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		view := tea.NewView("Loading...")
		view.AltScreen = true
		return view
	}

	bodyHeight := max(minBodyHeight, m.height-headerHeight-footerHeight)
	leftWidth := max(minOutputPaneWidth, m.width*outputPanePercent/percentScale)
	rightWidth := max(minCommentsPaneWidth, m.width-leftWidth-paneGapWidth)

	header := headerStyle.Width(m.width).Render(fmt.Sprintf("rvw: %s  exit=%d", m.command, m.exitCode))
	left := m.renderOutput(leftWidth, bodyHeight)
	right := m.renderComments(rightWidth, bodyHeight)
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := footerStyle.Width(m.width).Render(m.footerText())

	view := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, header, body, footer))
	view.AltScreen = true
	return view
}

func (m Model) renderOutput(width int, height int) string {
	lines := make([]string, 0, height)
	end := min(len(m.lines), m.offset+height)
	for i := m.offset; i < end; i++ {
		line := m.lines[i]
		marker := emptyCommentMarker
		if _, ok := m.comments.Get(line.Number); ok {
			marker = lineCommentMarker
		}
		text := truncate(fmt.Sprintf("%s%*d %s", marker, lineNumberWidth, line.Number, line.Plain), width)
		if i == m.cursor {
			text = currentLineStyle.Width(width).Render(text)
		} else {
			text = outputStyle.Width(width).Render(text)
		}
		lines = append(lines, text)
	}
	for len(lines) < height {
		lines = append(lines, outputStyle.Width(width).Render(""))
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderComments(width int, height int) string {
	items := m.comments.All()
	lines := make([]string, 0, height)
	lines = append(lines, sidebarTitleStyle.Width(width).Render("Comments"))
	for _, comment := range items {
		if len(lines) >= height {
			break
		}
		text := fmt.Sprintf("L%d: %s", comment.Line, comment.Text)
		lines = append(lines, sidebarStyle.Width(width).Render(truncate(text, width)))
	}
	for len(lines) < height {
		lines = append(lines, sidebarStyle.Width(width).Render(""))
	}
	return strings.Join(lines, "\n")
}

func (m Model) footerText() string {
	if m.editing {
		return "comment> " + m.input
	}
	if m.status != "" {
		return m.status + " | " + defaultFooterHelp
	}
	return defaultFooterHelp
}

func (m Model) keepCursorVisible() Model {
	bodyHeight := max(minBodyHeight, m.height-headerHeight-footerHeight)
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+bodyHeight {
		m.offset = m.cursor - bodyHeight + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
	return m
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}
	if width == 1 {
		return string(runes[:1])
	}
	return string(runes[:width-1]) + "…"
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true)
	footerStyle = lipgloss.NewStyle().
			Faint(true)
	outputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
	currentLineStyle = lipgloss.NewStyle().
				Reverse(true)
	sidebarTitleStyle = lipgloss.NewStyle().
				Bold(true)
	sidebarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))
)
