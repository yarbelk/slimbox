package wc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/text/width"
)

const order = "lwmcL"

const (
	LineFlag = "lines"
	ByteFlag = "bytes"
	WordFlag = "words"
	CharFlag = "chars"
	LongFlag = "max-line-length"
)

type Options interface {
	Args() []string
	GetBool(string) (bool, error)
	NFlag() int
}

type DefaultOptions struct{ Filenames []string }

func (d DefaultOptions) NFlag() int {
	return 0
}

func (d DefaultOptions) GetBool(n string) (bool, error) {
	switch n {
	case LineFlag, ByteFlag, WordFlag:
		return true, nil
	default:
		return false, nil
	}
}

func (d DefaultOptions) Args() []string {
	if d.Filenames == nil {
		return []string{}
	}
	return d.Filenames
}

// Results are the totals for output
type Results struct {
	Bytes, Characters, Newlines, Words, Longest uint

	Filename string
}

type ResultsSet struct {
	MaxNumber uint

	Results []Results
}

func approxLog10(u uint) uint {
	var i uint = 1
	for ; true; i++ {
		u /= 10
		if u == 0 {
			return i
		}
	}
	return 1
}

func (r Results) GetMax(options Options) (max uint) {
	formatOpts := options
	if options.NFlag() == 0 {
		formatOpts = DefaultOptions{}
	}
	for _, c := range order[:4] {
		switch c {
		case 'l':
			if ok, err := formatOpts.GetBool(LineFlag); !ok || err != nil {
				break
			}
			if max < r.Newlines {
				max = r.Newlines
			}
		case 'w':
			if ok, err := formatOpts.GetBool(WordFlag); !ok || err != nil {
				break
			}
			if max < r.Words {
				max = r.Words
			}
		case 'm':
			if ok, err := formatOpts.GetBool(CharFlag); !ok || err != nil {
				break
			}
			if max < r.Characters {
				max = r.Characters
			}
		case 'c':
			if ok, err := formatOpts.GetBool(ByteFlag); !ok || err != nil {
				break
			}
			if max < r.Bytes {
				max = r.Bytes
			}
		}
	}
	return
}

// Printf prints all results based on format options
func (rs ResultsSet) Printf(options Options) string {
	builder := strings.Builder{}
	width := approxLog10(rs.MaxNumber) + 2
	fmtstring := fmt.Sprintf("%%%dd", width)
	formatOpts := options
	if options.NFlag() == 0 {
		formatOpts = DefaultOptions{}
	}
	for _, results := range rs.Results {
	loop:
		for _, c := range order {
			switch c {
			case 'l':
				if ok, err := formatOpts.GetBool(LineFlag); !ok || err != nil {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Newlines))
			case 'w':
				if ok, err := formatOpts.GetBool(WordFlag); !ok || err != nil {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Words))
			case 'm':
				if ok, err := formatOpts.GetBool(CharFlag); !ok || err != nil {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Characters))
			case 'c':
				if ok, err := formatOpts.GetBool(ByteFlag); !ok || err != nil {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Bytes))
			case 'L':
				if ok, err := formatOpts.GetBool(LongFlag); !ok || err != nil {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Longest))
			}

		}
		if results.Filename != "" {
			builder.WriteString("  ")
			builder.WriteString(results.Filename)
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

// Default layout
func (rs ResultsSet) String() string {
	builder := strings.Builder{}
	width := approxLog10(rs.MaxNumber) + 2
	fmtstring := fmt.Sprintf("%%%dd%%%dd%%%dd  %%s", width, width, width)
	for _, result := range rs.Results {
		builder.WriteString(fmt.Sprintf(fmtstring, result.Newlines, result.Words, result.Bytes, result.Filename))
		builder.WriteRune('\n')
	}
	return builder.String()
}

func runeSize(r rune) uint {
	switch width.LookupRune(r).Kind() {
	case width.EastAsianNarrow:
		return 1
	case width.EastAsianFullwidth, width.EastAsianWide:
		return 2
	default:
		return 1
	}
}

func WordCount(opts Options, in io.Reader) (Results, error) {
	buffered := bufio.NewReader(in)
	results := Results{}
	var inWord, position uint

	for {
		r, size, err := buffered.ReadRune()
		if err != nil {
			if err == io.EOF {
				results.Words += inWord
				if position > results.Longest {
					results.Longest = position
				}
				return results, nil
			}
			return results, err
		}
		if size > 0 {
			results.Characters++
			results.Bytes += uint(size)
		}
		if unicode.IsPrint(r) && !unicode.IsSpace(r) {
			inWord = 1
			position += runeSize(r)
			continue
		}
		if unicode.IsSpace(r) {
			results.Words += inWord
			inWord = 0
			switch r {
			case '\t':
				// round up to 8 ( set LSBs to 111, then add one)
				position = (position | 7) + 1
			case '\v':
				results.Newlines++
			case '\u00A0', ' ':
				position += runeSize(r)
				continue
			case '\n':
				results.Newlines++
				if position > results.Longest {
					results.Longest = position
				}
				position = 0
			case '\r', '\u0085':
				if position > results.Longest {
					results.Longest = position
				}
				position = 0
			case '\f':
				if position > results.Longest {
					results.Longest = position
				}
				position = 0
			}
		}
	}
}
