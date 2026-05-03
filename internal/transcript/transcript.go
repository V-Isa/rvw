package transcript

import (
	"regexp"
	"strings"
)

type Line struct {
	Number int
	Raw    string
	Plain  string
}

type Transcript struct {
	Lines []Line
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

func FromBytes(data []byte) Transcript {
	raw := strings.ReplaceAll(string(data), "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.TrimSuffix(raw, "\n")
	if raw == "" {
		return Transcript{}
	}

	parts := strings.Split(raw, "\n")
	lines := make([]Line, 0, len(parts))
	for i, part := range parts {
		lines = append(lines, Line{
			Number: i + 1,
			Raw:    part,
			Plain:  StripANSI(part),
		})
	}
	return Transcript{Lines: lines}
}

func StripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}
