package goscrobble

import (
	"testing"

	"git.m2.nz/go-scrobble/internal/goscrobble"
)

func TestParseJellyfinInput(t *testing.T) {
	got := goscrobble.ParseJellyfinInput("TestString!")
	if got != "TestString!" {
		t.Errorf("ParseJellyfinInput returned: %s", got)
	}
}
