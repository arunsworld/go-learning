package learning

import (
	"path"
	"testing"
)

func TestFilenameFromPath(t *testing.T) {
	fullFileName := "/tmp/a/b/c.txt"

	filename := path.Base(fullFileName)
	expected := "c.txt"

	if filename != expected {
		t.Error("Expected c.txt, got:", filename)
	}
}
