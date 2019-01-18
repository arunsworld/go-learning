package learning

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCreatingAZipFile(t *testing.T) {

	type zipFile struct {
		filename string
		contents io.Reader
	}

	input := []zipFile{
		zipFile{filename: "file1.txt", contents: strings.NewReader("contents of file1.txt")},
		zipFile{filename: "file2.txt", contents: strings.NewReader("contents of file2.txt")},
		zipFile{filename: "file3.txt", contents: strings.NewReader("contents of file3.txt")},
	}

	var zipFileContents io.Writer
	zipFileContents = new(bytes.Buffer)
	// zipFileContents, _ = os.Create("a.zip")
	// defer zipFileContents.(*os.File).Close()

	zipWriter := zip.NewWriter(zipFileContents)
	for _, zf := range input {
		w, err := zipWriter.Create(zf.filename)
		if err != nil {
			t.Error("Errour in creating zip file entry: ", err)
			return
		}
		io.Copy(w, zf.contents)
	}
	zipWriter.Close()

	checkSum := fmt.Sprintf("%x", md5.Sum(zipFileContents.(*bytes.Buffer).Bytes()))
	if checkSum != "dc23b297c9b68e4cda841427c7ce9415" {
		t.Errorf("Checksum didn't match. Expecting %s. Got %s.", "dc23b297c9b68e4cda841427c7ce9415", checkSum)
	}

}

func TestReadingAZipFile(t *testing.T) {

	r, err := zip.OpenReader("zip_test.zip")
	if err != nil {
		t.Error("Error opening zip file: ", err)
		return
	}
	defer r.Close()

	for i, f := range r.File {
		switch i {
		case 0:
			if f.Name != "file1.txt" {
				t.Error("Expecting file1.txt, got: ", f.Name)
				return
			}
			ff, err := f.Open()
			if err != nil {
				t.Error("Could not open file1.txt for reading: ", err)
				return
			}
			byteContents, err := ioutil.ReadAll(ff)
			if err != nil {
				t.Error("Could not read contents of file1.txt: ", err)
				ff.Close()
				return
			}
			contents := string(byteContents)
			if contents != "contents of file1.txt" {
				t.Error("Incorrect contents of file1.txt. Got: ", contents)
			}
			ff.Close()
		case 1:
			if f.Name != "file2.txt" {
				t.Error("Expecting file2.txt, got: ", f.Name)
				return
			}
		case 2:
			if f.Name != "file3.txt" {
				t.Error("Expecting file3.txt, got: ", f.Name)
				return
			}
		}
	}

}

func TestReadingABadZipFile(t *testing.T) {

	_, err := zip.OpenReader("doesnotexist.zip")
	if err == nil {
		t.Error("Did not get an error when trying to open doesnotexist.zip.")
		return
	}
	if err.Error() != "open doesnotexist.zip: no such file or directory" {
		t.Error("Did not get the right error message when trying to open doesnotexist.zip")
		return
	}

	w, err := zip.OpenReader("zip_test.go")
	if w != nil {
		t.Error("Expecting nil, got: ", w)
	}
	if err == nil {
		t.Error("Did not get an error when opening bad zip file.")
		return
	}
	if err.Error() != "zip: not a valid zip file" {
		t.Error("Expecting zip: not a valid zip file. Got: ", err.Error())
	}
	if err != zip.ErrFormat {
		t.Error("Did not get the expected error variable.")
	}

}
