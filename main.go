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
	var err error

	file, err = filepath.Abs(flag.Arg(0))
	print_err(err)
	return file

}

func init() {
	flag.Parse()
	if flag.NArg() > 1 {
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Concatenate files and print them to stdout\n")
	fmt.Fprintf(os.Stderr, "\tUSAGE: %s [FILENAME]\n", os.Args[0])
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
		file = parse_file_arg(file)
		func() {
			fi, err := os.Open(file)
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
