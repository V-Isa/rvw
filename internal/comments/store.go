package comments

import (
	"sort"
	"strings"
)

type Comment struct {
	Line int
	Text string
}

type Store struct {
	byLine map[int]string
}

func NewStore() *Store {
	return &Store{byLine: make(map[int]string)}
}

func (s *Store) Set(line int, text string) {
	text = strings.TrimSpace(text)
	if line < 1 {
		return
	}
	if text == "" {
		delete(s.byLine, line)
		return
	}
	s.byLine[line] = text
}

func (s *Store) Get(line int) (string, bool) {
	text, ok := s.byLine[line]
	return text, ok
}

func (s *Store) All() []Comment {
	lines := make([]int, 0, len(s.byLine))
	for line := range s.byLine {
		lines = append(lines, line)
	}
	sort.Ints(lines)

	out := make([]Comment, 0, len(lines))
	for _, line := range lines {
		out = append(out, Comment{Line: line, Text: s.byLine[line]})
	}
	return out
}
