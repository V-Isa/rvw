package transcript

import "testing"

func TestFromBytesStoresRawAndPlainLines(t *testing.T) {
	tr := FromBytes([]byte("\x1b[31mError\x1b[0m\r\nok\n"))
	if len(tr.Lines) != 2 {
		t.Fatalf("len = %d, want 2", len(tr.Lines))
	}
	if tr.Lines[0].Raw != "\x1b[31mError\x1b[0m" {
		t.Fatalf("raw = %q, want %q", tr.Lines[0].Raw, "\x1b[31mError\x1b[0m")
	}
	if tr.Lines[0].Plain != "Error" {
		t.Fatalf("plain = %q, want %q", tr.Lines[0].Plain, "Error")
	}
	if tr.Lines[1].Number != 2 || tr.Lines[1].Plain != "ok" {
		t.Fatalf("second line = %#v, want number = 2, plain = %q", tr.Lines[1], "ok")
	}
}
