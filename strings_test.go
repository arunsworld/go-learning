package learning

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestStringFormatting(t *testing.T) {
	type point struct {
		x int
		y int
	}
	p := point{x: 3, y: 5}

	t.Run("Structs", func(t *testing.T) {
		got := fmt.Sprintf("%v", p)
		expected := "{3 5}"
		validate(t, got, expected)

		got = fmt.Sprintf("%+v", p)
		expected = "{x:3 y:5}"
		validate(t, got, expected)

		got = fmt.Sprintf("%#v", p)
		expected = "learning.point{x:3, y:5}"
		validate(t, got, expected)
	})

	t.Run("Pointers", func(t *testing.T) {
		t.Skip("Documentation only.")
		got := fmt.Sprintf("%p", &p)
		expected := "0xc000014100"
		validate(t, got, expected)
	})

	t.Run("TypeIntrospection", func(t *testing.T) {
		got := fmt.Sprintf("%T", p)
		expected := "learning.point"
		validate(t, got, expected)
	})

	t.Run("Numbers", func(t *testing.T) {
		got := fmt.Sprintf("%x", 255)
		expected := "ff"
		validate(t, got, expected)

		got = fmt.Sprintf("%X", 255)
		expected = "FF"
		validate(t, got, expected)

		got = fmt.Sprintf("%x", 0xff)
		expected = "ff"
		validate(t, got, expected)

		got = fmt.Sprintf("%0.2f", 123.1294)
		expected = "123.13"
		validate(t, got, expected)
	})

	t.Run("Unicode", func(t *testing.T) {
		got := fmt.Sprintf("%+q", "✓")
		expected := `"\u2713"`
		validate(t, got, expected)

		got = fmt.Sprintf("%U", '✓')
		expected = "U+2713"
		validate(t, got, expected)

		r, _ := utf8.DecodeRune([]byte{0xe2, 0x9c, 0x93})
		got = fmt.Sprintf("%c", r)
		expected = "✓"
		validate(t, got, expected)
	})

}

func validate(t *testing.T, got, expected string) {
	if got != expected {
		t.Fatalf("Exected %s but got: %s", expected, got)
	}
}
