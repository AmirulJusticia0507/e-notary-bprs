package hash

import (
	"strings"
	"testing"
)

func TestFromString(t *testing.T) {
	got := FromString("hello")
	// SHA-256("hello") yang sudah dikenal.
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != want {
		t.Errorf("FromString() = %q, want %q", got, want)
	}
	if len(got) != 64 {
		t.Errorf("hash length = %d, want 64", len(got))
	}
}

func TestFromReader(t *testing.T) {
	got, err := FromReader(strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("FromReader() error = %v", err)
	}
	if want := FromString("hello"); got != want {
		t.Errorf("FromReader() = %q, want %q", got, want)
	}
}

func TestFromBytes(t *testing.T) {
	if got, want := FromBytes([]byte("hello")), FromString("hello"); got != want {
		t.Errorf("FromBytes() = %q, want %q", got, want)
	}
}
