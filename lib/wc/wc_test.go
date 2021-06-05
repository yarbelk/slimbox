package wc_test

import (
	"strings"
	"testing"

	"gitlab.com/yarbelk/slimbox/lib/wc"
)

const (
	arabic  = "\u0645\u0631\u062d\u0628\u0627\u0020\u0628\u0627\u0644\u0639\u0627\u0644\u0645\u0021"
	chinese = "\u4f60\u597d\uff0c\u4e16\u754c\uff01"
)

type TestOptions struct {
	flags map[string]bool
	args  []string
}

func (to TestOptions) GetBool(name string) (bool, error) {
	return to.flags[name], nil
}

func (to TestOptions) Args() []string {
	return to.args
}

func TestBasicTypesStdinStyle(t *testing.T) {
	var tests = []struct {
		name     string
		expected wc.Results
		given    *strings.Reader
		args     []string
	}{
		{"Empty File", wc.Results{}, strings.NewReader(""), []string{}},
		{"Once Char: 1 c/b/w/l", wc.Results{Bytes: 1, Characters: 1, Words: 1, Newlines: 0, Longest: 1}, strings.NewReader("a"), []string{}},
		{"Two Char: 2 c/b 1 w/l", wc.Results{Bytes: 2, Characters: 2, Words: 1, Newlines: 0, Longest: 2}, strings.NewReader("ab"), []string{}},
		{"Three Char: 3 c/b 1 w/l", wc.Results{Bytes: 3, Characters: 3, Words: 1, Newlines: 0, Longest: 3}, strings.NewReader("abc"), []string{}},
		{"Two Words: 3 c/b 2 w 1 l", wc.Results{Bytes: 3, Characters: 3, Words: 2, Newlines: 0, Longest: 3}, strings.NewReader("a b"), []string{}},
		{"words on Lines: 3 c/b 2 w 2 l", wc.Results{Bytes: 3, Characters: 3, Words: 2, Newlines: 1, Longest: 1}, strings.NewReader("a\nb"), []string{}},
		{"Tabs for Length: 3 c/b 2 w 2 l 1", wc.Results{Bytes: 3, Characters: 3, Words: 2, Newlines: 0, Longest: 9}, strings.NewReader("a\tb"), []string{}},
		{"Multiple Tabs for Length: 3 c/b 2 w 2 l 2", wc.Results{Bytes: 4, Characters: 4, Words: 2, Newlines: 0, Longest: 17}, strings.NewReader("a\t\tb"), []string{}},
		{"Vertical Tab for confusion: 3 c/b 2 w 2 l 2", wc.Results{Bytes: 3, Characters: 3, Words: 2, Newlines: 1, Longest: 2}, strings.NewReader("a\vb"), []string{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			options := TestOptions{make(map[string]bool), tt.args}
			actual, err := wc.WordCount(options, tt.given)
			if err != nil {
				t.Errorf("expected no error, actual %s", err)
			}
			if actual != tt.expected {
				t.Errorf("expected %+v, actual %+v", tt.expected, actual)
			}
		})
	}
}

func TestUnicodeStuff(t *testing.T) {
	var tests = []struct {
		name     string
		expected wc.Results
		given    *strings.Reader
		args     []string
	}{
		{"arabic (note: no idea the right way, just duplicate gnu wc", wc.Results{Bytes: 26, Characters: 14, Words: 2, Newlines: 0, Longest: 14}, strings.NewReader(arabic), []string{}},
		{"Chinese", wc.Results{Bytes: 18, Characters: 6, Words: 1, Newlines: 0, Longest: 12}, strings.NewReader(chinese), []string{}},
		{"zero width space", wc.Results{Bytes: 5, Characters: 3, Words: 1, Newlines: 0, Longest: 2}, strings.NewReader("a\u200Bb"), []string{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Log(arabic)
			t.Log(chinese)
			options := TestOptions{make(map[string]bool), tt.args}
			actual, err := wc.WordCount(options, tt.given)
			if err != nil {
				t.Errorf("expected no error, actual %s", err)
			}
			if actual != tt.expected {
				t.Errorf("\n\t\texpected %+v\n\t\tactual   %+v", tt.expected, actual)
			}
		})
	}
}
