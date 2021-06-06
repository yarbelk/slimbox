package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"gitlab.com/yarbelk/slimbox/lib/wc"
)

// I don't like the use of the wc.XFlag constants here: its too coupled
var (
	options      *pflag.FlagSet = pflag.NewFlagSet("wc", pflag.ContinueOnError)
	lines        *bool          = options.BoolP(wc.LineFlag, "l", false, "Count newlines")
	bytes        *bool          = options.BoolP(wc.ByteFlag, "c", false, "Count bytes")
	words        *bool          = options.BoolP(wc.WordFlag, "w", false, "Count words")
	characters   *bool          = options.BoolP(wc.CharFlag, "m", false, "Count characters")
	printLongist *bool          = options.BoolP(wc.LongFlag, "L", false, "Print longest line length")
)

func main() {
	if err := options.Parse(os.Args[1:]); err != nil {
		fmt.Fprint(os.Stderr, err)
		options.Usage()
	}
	if err := wc.WcMain(options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
