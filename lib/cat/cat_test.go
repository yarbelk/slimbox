package cat_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yarbelk/slimbox/lib/cat"
)

func TestOneFileInputs(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    *strings.Reader
	}{
		{"Empty File", "", strings.NewReader("")},
		{"One Line File with newline", "Hello world!\n", strings.NewReader("Hello world!\n")},
		{"One Line File with no trailing newline", "Hello world!\n", strings.NewReader("Hello world!")},
		{"Multiple lines with trailing newline", "Hello world!\n\nHello Gophers!\n", strings.NewReader("Hello world!\n\nHello Gophers!\n")},
		{"Multiple lines with no trailing newline", "Hello world!\n\nHello Gophers!\n", strings.NewReader("Hello world!\n\nHello Gophers!")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			outputStream := &bytes.Buffer{}

			C := &cat.CatOptions{}

			C.Cat(tt.given, outputStream)
			actual := outputStream.String()
			if actual != tt.expected {
				t.Errorf("expected %s, actual %s", tt.expected, actual)
			}
		})
	}
}

func TestMultipleInputs(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    []*strings.Reader
	}{
		{"Two Empty Files", "", []*strings.Reader{strings.NewReader(""), strings.NewReader("")}},
		{"One Line Files with newline", "Hello world!\nHello Gophers!\n", []*strings.Reader{strings.NewReader("Hello world!\n"), strings.NewReader("Hello Gophers!\n")}},
		{"One Line Files with no trailing newline", "Hello world!\nHello Gophers!\n", []*strings.Reader{strings.NewReader("Hello world!"), strings.NewReader("Hello Gophers!")}},
		{"Mixed trailing newlines", "Hello world!\nHello Gophers!\n", []*strings.Reader{strings.NewReader("Hello world!\n"), strings.NewReader("Hello Gophers!")}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			outputStream := &bytes.Buffer{}

			C := &cat.CatOptions{}

			for _, input := range tt.given {
				C.Cat(input, outputStream)
			}
			actual := outputStream.String()
			if actual != tt.expected {
				t.Errorf("expected %s, actual %s", tt.expected, actual)
			}
		})
	}
}
