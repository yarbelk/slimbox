package wc_test

import (
	"testing"

	"gitlab.com/yarbelk/slimbox/lib/wc"
)

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
			Options:    wc.Options{Files: []string{"testfile"}},
			ResultsSet: resultset,
		}, {
			name:       "line and words",
			expected:   "   2   2  testfile\n",
			Options:    wc.Options{Files: []string{"testfile"}, Newlines: true, Words: true, NFlag: 2},
			ResultsSet: resultset,
		}, {
			name:       "two lines, scaled correctly",
			expected:   "    2    2    5   10    3  testfile\n   22   18   50  900    8  testfile2\n",
			Options:    wc.Options{Files: []string{"testfile", "testfile2"}, Bytes: true, Characters: true, Longest: true, Newlines: true, Words: true, NFlag: 4},
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
