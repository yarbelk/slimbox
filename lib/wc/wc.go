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

type Options struct {
	Bytes, Characters, Newlines, Words, Longest bool

	NFlag uint
	Files []string
}

func (wo Options) GetBool(flag string) bool {
	switch flag {
	case ByteFlag:
		return wo.NFlag == 0 || wo.Bytes
	case CharFlag:
		return wo.Characters
	case WordFlag:
		return wo.NFlag == 0 || wo.Words
	case LineFlag:
		return wo.NFlag == 0 || wo.Newlines
	case LongFlag:
		return wo.Longest
	default:
		return false
	}
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

// GetMax number (that isn't line) from this resultset for use in max line length parameter
func (r Results) GetMax(options Options) (max uint) {
	for _, c := range order[:4] {
		switch c {
		case 'l':
			if !options.GetBool(LineFlag) {
				break
			}
			if max < r.Newlines {
				max = r.Newlines
			}
		case 'w':
			if !options.GetBool(WordFlag) {
				break
			}
			if max < r.Words {
				max = r.Words
			}
		case 'm':
			if !options.GetBool(CharFlag) {
				break
			}
			if max < r.Characters {
				max = r.Characters
			}
		case 'c':
			if !options.GetBool(ByteFlag) {
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
	for _, results := range rs.Results {
	loop:
		for _, c := range order {
			switch c {
			case 'l':
				if !options.GetBool(LineFlag) {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Newlines))
			case 'w':
				if !options.GetBool(WordFlag) {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Words))
			case 'm':
				if !options.GetBool(CharFlag) {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Characters))
			case 'c':
				if !options.GetBool(ByteFlag) {
					continue loop
				}
				builder.WriteString(fmt.Sprintf(fmtstring, results.Bytes))
			case 'L':
				if !options.GetBool(LongFlag) {
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
