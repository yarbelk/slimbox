package main

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	"io"
	"os"
	"path/filepath"
)

func print_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		usage()
	}
}

func parse_file_arg(file string) string {
	full_file, err := filepath.Abs(file)
	print_err(err)
	return full_file
}

func init() {
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Concatenate file[s] or standard input to stdout\n")
	fmt.Fprintf(os.Stderr, "\tUSAGE: %s [FILENAME]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "With no file, or when file is `-', read from stdin", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
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
		nn, err := os.Stdout.Write(line)

		if nn < len(line) {
			print_err(err)
		}
	}
	return nil
}

func main() {
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
