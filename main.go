package main

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	"io"
	"os"
	"path/filepath"
)

var file string

var output *bufio.Writer

func print_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		usage()
	}
}

func parse_file_arg() string {
	var err error

	file, err = filepath.Abs(flag.Arg(0))
	print_err(err)
	return file

}

func init() {
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Concatenate files and print them to stdout\n")
	fmt.Fprintf(os.Stderr, "\tUSAGE: %s [FILENAME]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func cat_file(filename string) error {
	var err error

	fi, err := os.Open(file)
	print_err(err)

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	input_buffer := bufio.NewReader(fi)

	for line, err := input_buffer.ReadBytes(byte('\n')); err != nil; {
		nn, err := output.Write(line)

		if nn < len(line) {
			print_err(err)
		}
	}
	if err != io.EOF {
		print_err(err)
	}
	return nil
}

func main() {
	output = bufio.NewWriter(os.Stdout)
	parse_file_arg()

	_ = cat_file(file)
}
