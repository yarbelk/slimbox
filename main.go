package main

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	"io"
	"os"
	"path/filepath"
)

var numberLines, help *bool
var lineNumber int = 1

func usage() {
	fmt.Fprintf(os.Stderr, "Concatenate file[s] or standard input to standard output\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "With no FILE, or when FILE is -, read standard input.\n")

	fmt.Fprintf(os.Stderr, "\nExamples:\n\n")
	fmt.Fprintf(os.Stderr, "    cat f - g  Output f's contents, then standard input, then g's contents.\n")
	fmt.Fprintf(os.Stderr, "    cat        Copy standard input to standard output.\n")
	os.Exit(2)
}

func print_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}

func parse_file_arg(file string) string {
	full_file, err := filepath.Abs(file)
	print_err(err)
	return full_file
}

func init() {
	numberLines = flag.BoolP("number", "n", false, "number all output lines")
	help = flag.BoolP("help", "h", false, "print this message")

	flag.Parse()
}

func cat_file(fi *os.File) error {
	input_buffer := bufio.NewReader(fi)

	for {
		line, err := input_buffer.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				print_err(err)
			}
			return nil
		}
		if *numberLines {
			line = append([]byte(fmt.Sprintf("%6d  ", lineNumber)), line...)
		}
		nn, err := os.Stdout.Write(line)

		if nn < len(line) {
			print_err(err)
		}
		lineNumber += 1
	}
	return nil
}

func main() {
	if *help {
		usage()
	}
	if flag.NArg() == 0 {
		cat_file(os.Stdin)
	}

	for _, file := range flag.Args() {
		if file == "-" {
			_ = cat_file(os.Stdin)
		} else {
			func() {
				full_file := parse_file_arg(file)
				fi, err := os.Open(full_file)
				print_err(err)

				_ = cat_file(fi)

				defer func() {
					if err := fi.Close(); err != nil {
						print_err(err)
					}
				}()
			}()
		}

	}
}
