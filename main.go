package main

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	"github.com/yarbelk/cat/lib"
	"os"
)

var conf lib.Conf

var numberLines, numberNonBlankLines, help bool
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

func init() {
	conf = lib.Conf{false, false, false, 1}
	flag.BoolVarP(&conf.NumberLines, "number", "n", false, "number all output lines")
	flag.BoolVarP(&conf.NumberNonBlankLines, "number-nonblank", "b", false, "number non-blank output lines, overrides -n")
	flag.BoolVarP(&conf.Help, "help", "h", false, "print this message")

	flag.Parse()

	if conf.NumberLines && conf.NumberNonBlankLines {
		conf.NumberLines = false
	}
}

func main() {
	if help {
		usage()
	}
	if flag.NArg() == 0 {
		conf.CatFile(bufio.NewReader(os.Stdin))
	}

	for _, file := range flag.Args() {
		func() {
			input_buffer, err := lib.ParseFileArg(file)
			if err != nil {
				lib.PrintErr(err)
				os.Exit(1)
			}

			_ = conf.CatFile(input_buffer)

		}()

	}
}
