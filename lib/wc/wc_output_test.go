package wc_test

import (
	"testing"

	"gitlab.com/yarbelk/slimbox/lib/wc"
)

type testOptions struct {
	Opts      map[string]bool
	nset      int
	filenames []string
}

func (t *testOptions) NFlag() int {
	return t.nset
}

func (t *testOptions) GetBool(n string) (bool, error) {
	return t.Opts[n], nil
}

func (t *testOptions) Args() []string {
	return t.filenames
}

func TestDefaultResultFormating(t *testing.T) {
	resultset := wc.ResultsSet{MaxNumber: 10, Results: []wc.Results{{Bytes: 10, Characters: 5, Newlines: 2, Words: 2, Longest: 3, Filename: "testfile"}}}
	expected := "   2   2  10  testfile\n"
	if resultset.String() != expected {
		t.Errorf("\n\t\texpected %+v\n\t\tactual   %+v", expected, resultset.String())
	}
}

func TestFormating(t *testing.T) {
	resultset := wc.ResultsSet{MaxNumber: 10, Results: []wc.Results{{Bytes: 10, Characters: 5, Newlines: 2, Words: 2, Longest: 3, Filename: "testfile"}}}
	resultset2 := wc.ResultsSet{MaxNumber: 900, Results: []wc.Results{
		{Bytes: 10, Characters: 5, Newlines: 2, Words: 2, Longest: 3, Filename: "testfile"},
		{Bytes: 900, Characters: 50, Newlines: 22, Words: 18, Longest: 8, Filename: "testfile2"},
	},
	}
	var tests = []struct {
		name       string
		expected   string
		Options    wc.Options
		ResultsSet wc.ResultsSet
	}{
		{
			name:       "Standard",
			expected:   "   2   2  10  testfile\n",
			Options:    &testOptions{Opts: make(map[string]bool)},
			ResultsSet: resultset,
		}, {
			name:       "line and words",
			expected:   "   2   2  testfile\n",
			Options:    &testOptions{Opts: map[string]bool{"l": true, "w": true}, nset: 2},
			ResultsSet: resultset,
		}, {
			name:       "two lines, scaled correctly",
			expected:   "    2    2    5   10    3  testfile\n   22   18   50  900    8  testfile2\n",
			Options:    &testOptions{Opts: map[string]bool{"l": true, "w": true, "c": true, "m": true, "L": true}, nset: 2},
			ResultsSet: resultset2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.ResultsSet.Printf(tt.Options)
			if actual != tt.expected {
				t.Errorf("given(%+v):\n\t\texpected\n%s\n\t\tactual:\n %s", tt.Options, tt.expected, actual)
			}
		})
	}
}
