package comments

import "testing"

func TestStoreSetGetDelete(t *testing.T) {
	store := NewStore()

	store.Set(42, " check this ")
	got, ok := store.Get(42)
	if !ok {
		t.Fatal("comment exists = false, want true")
	}
	if got != "check this" {
		t.Fatalf("comment = %q, want %q", got, "check this")
	}

	store.Set(42, "")
	if _, ok := store.Get(42); ok {
		t.Fatal("comment exists = true, want false")
	}
}

func TestStoreAllSorted(t *testing.T) {
	store := NewStore()
	store.Set(7, "second")
	store.Set(2, "first")

	got := store.All()
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Line != 2 || got[1].Line != 7 {
		t.Fatalf("comments = %#v, want sorted by line", got)
	}
}
