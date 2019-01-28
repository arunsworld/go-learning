package learning

import (
	"bytes"
	"fmt"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestUTF8Encodng(t *testing.T) {
	var got []byte
	var expected []byte

	p := make([]byte, 6)
	count := utf8.EncodeRune(p, '✓')

	got = p[:count]
	expected = []byte{0xe2, 0x9c, 0x93}
	validateByByteComparison(t, got, expected)

	got = p[:count]
	expected = []byte("\xe2\x9c\x93")
	validateByByteComparison(t, got, expected)

	gotS := fmt.Sprintf("%x", p[:count])
	expectedS := "e29c93"
	validateString(t, gotS, expectedS)

	gotS = fmt.Sprintf("%s", p[:count])
	expectedS = "✓"
	validateString(t, gotS, expectedS)

	gotS = fmt.Sprintf("%+q", p[:count])
	expectedS = `"\u2713"`
	validateString(t, gotS, expectedS)

}

func TestUnicodeTests(t *testing.T) {

	a := '✓'
	fmt.Printf("(%c) Lower: (%c)\n", a, unicode.ToUpper(a))

}

func validateByByteComparison(t *testing.T, got, expected []byte) {
	if bytes.Compare(got, expected) != 0 {
		t.Fatalf("Expected %X, got %X", expected, got)
	}
}

func validateString(t *testing.T, got, expected string) {
	if got != expected {
		t.Fatalf("Exected %s but got: %s", expected, got)
	}
}
