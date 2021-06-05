package wc

import (
	"bufio"
	"io"
	"unicode"

	"golang.org/x/text/width"
)

type Options interface {
	Args() []string
	GetBool(name string) (bool, error)
}

// Results are the totals for output
type Results struct {
	Bytes, Characters, Newlines, Words, Longest uint
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
	return results, nil
}
