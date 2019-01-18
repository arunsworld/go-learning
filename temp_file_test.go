package learning

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreatingATempFile(t *testing.T) {
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("", "example.txt")
	if err != nil {
		t.Fatal("Unable to create temp file:", err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal("Unable to write to temp file:", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Unable to close temp file:", err)
	}
}
