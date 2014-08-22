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
func ParseFileArg(file string) (*bufio.Reader, error) {
	var input_buffer *bufio.Reader
	if file == "-" {
		return bufio.NewReader(os.Stdin), nil
	}

	full_file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	fi, err := os.Open(full_file)
	if err != nil {
		return nil, err
	}

	input_buffer = bufio.NewReader(fi)
	return input_buffer, nil
}

func (c *Conf) CatFile(input_buffer *bufio.Reader) error {
	var number []byte
	for {
		line, err := input_buffer.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				PrintErr(err)
			}
			return nil
		}

		if c.NumberLines {
			line = append([]byte(fmt.Sprintf("%6d  ", c.LineNumber)), line...)
			c.LineNumber += 1
		} else if c.NumberNonBlankLines {

			if len(strings.TrimSpace(string(line))) != 0 {
				number = []byte(fmt.Sprintf("%6d  ", c.LineNumber))
				c.LineNumber += 1
			} else {
				number = []byte("      ")
			}
			line = append(number, line...)
		}
		_, err = os.Stdout.Write(line)
		if err != nil {
			return err
		}
	}
}
