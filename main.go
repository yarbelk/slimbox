package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"github.com/yarbelk/cat/lib"
	"log"
	"os"
)

var catOptions cat.CatOptions

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
	catOptions = cat.CatOptions{false, false}
	flag.BoolVarP(&catOptions.EoL, "show-ends", "E", false, "display $ at end of each line")
	flag.BoolVarP(&catOptions.Tabs, "show-tabs", "T", false, "display TAB characters as ^I")

	//flag.BoolVarP(&catOptions.NumberLines, "number", "n", false, "number all output lines")
	//flag.BoolVarP(&conf.NumberNonBlankLines, "number-nonblank", "b", false, "number non-blank output lines, overrides -n")
	flag.BoolVarP(&help, "help", "h", false, "print this message")

	flag.Parse()

	//if conf.NumberLines && conf.NumberNonBlankLines {
	//	conf.NumberLines = false
	//}
}

func main() {
	if help {
		usage()
	}
	if flag.NArg() == 0 {
		catOptions.Cat(os.Stdin, os.Stdout)
	}

	for _, file := range flag.Args() {
		func() {
			fi, err := os.Open(file)
			defer fi.Close()
			if err != nil {
				log.Fatal(err)
			}
			err = catOptions.Cat(fi, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		}()

	}
}
