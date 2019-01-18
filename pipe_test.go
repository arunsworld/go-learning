package learning

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func TestInputPipe(t *testing.T) {
}

func processPipeData() string {
	info, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return ""
	}
	if info.Size() == 0 {
		return ""
	}
	r := bufio.NewReader(os.Stdin)
	result, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	return string(result)
}
