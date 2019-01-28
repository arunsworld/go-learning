package learning

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func TestInputPipe(t *testing.T) {
	t.Skip("Don't know how to test so using this for documentation. Function below used to read data piped in.")
	processPipeData()
}

func processPipeData() string {
	// Sourced from http://blog.ralch.com/tutorial/golang-command-line-pipes/
	info, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return ""
	}
	r := bufio.NewReader(os.Stdin)
	result, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	return string(result)
}
