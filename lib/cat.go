package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Conf struct {
	NumberLines, NumberNonBlankLines, Help bool
	LineNumber                             int
}

func PrintErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}

// ParseFileArg takes in a string and returns a ByteReader object.
// if it is passed '-', it returns Stdin.
func ParseFileArg(file string) (*os.File, error) {
	if file == "-" {
		return os.Stdin, nil
	}

	full_file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	fi, err := os.Open(full_file)
	if err != nil {
		return nil, err
	}

	return fi, nil
}

func (c *Conf) numbers(line *[]byte) *[]byte {
	var number []byte
	var newline []byte
	if c.NumberLines {
		newline = append([]byte(fmt.Sprintf("%6d  ", c.LineNumber)), *line...)
		c.LineNumber += 1
	} else if c.NumberNonBlankLines {

		if len(strings.TrimSpace(string(*line))) != 0 {
			number = []byte(fmt.Sprintf("%6d  ", c.LineNumber))
			c.LineNumber += 1
		} else {
			number = []byte("      ")
		}
		newline = append(number, *line...)
	}
	return &newline
}

func (c *Conf) CatFile(fi *os.File) error {
	var input_buffer *bufio.Reader
	var err error = nil
	var line []byte
	input_buffer = bufio.NewReader(fi)
	for {

		line, err = input_buffer.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				PrintErr(err)
				return err
			}
			return nil
		}
		if c.NumberLines || c.NumberNonBlankLines {
			line = *c.numbers(&line)
		}

		_, err = os.Stdout.Write(line)
		if err != nil {
			return err
		}
	}
}
