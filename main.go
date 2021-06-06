package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/voxelbrain/goptions"
	"github.com/yarbelk/slimbox/lib"
	"github.com/yarbelk/slimbox/lib/cat"
	"github.com/yarbelk/slimbox/lib/falsy"
	"github.com/yarbelk/slimbox/lib/truthy"
	"gitlab.com/yarbelk/slimbox/lib/wc"
)

var catOptions cat.CatOptions

var numberLines, numberNonBlankLines, help bool
var lineNumber int = 1

type Options struct {
	Verb  goptions.Verbs
	Cat   cat.CatOptions     `goptions:"cat"`
	False falsy.FalseOptions `goptions:"false"`
	True  truthy.TrueOptions `goptions:"true"`
	Wc    struct{}           `goptions:"wc"`
}

func version() {
	os.Stderr.WriteString("Nice try - no versions yet\n")
	os.Exit(0)
}

func runFalse(options *Options) {
	if options.False.Version {
		version()
	}
	falsy.False()
}

func runTrue(options *Options) {
	if options.True.Version {
		version()
	}
	truthy.True()
}

func runCat(options *Options) {
	catOptions = *cat.NewCatOptions()
	catOptions.Blank = options.Cat.Blank
	catOptions.Number = options.Cat.Number
	catOptions.EoL = options.Cat.EoL
	catOptions.Tabs = options.Cat.Tabs

	if catOptions.Blank {
		catOptions.Number = true
	}
	if len(options.Cat.Files) == 0 {
		infile := os.NewFile(uintptr(syscall.Stdin), "/dev/stdin")
		err := catOptions.Cat(infile, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, file := range options.Cat.Files {
		func() {
			_, fi, err := lib.ParseFiles(file)
			if err != nil {
				log.Fatal(err)
			}
			defer fi.Close()
			err = catOptions.Cat(fi, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

}

func main() {
	options := Options{}

	flagset := goptions.NewFlagSet(filepath.Base(os.Args[0]), &options)
	err := flagset.Parse(os.Args[1:])

	// This is all placeholder until i migrate true, false and cat to pflag (fix argument ordering)
	switch options.Verb {
	case "cat":
		if err == nil {
			runCat(&options)
			os.Exit(0)
		}
		var catHelp goptions.HelpFunc = goptions.NewTemplatedHelpFunc(cat.CAT_HELP)
		flagset.Verbs["cat"].HelpFunc = catHelp
		flagset.Verbs["cat"].PrintHelp(os.Stderr)
		os.Exit(1)
	case "false":
		if err == nil {
			runFalse(&options)
		}
		var falseHelp goptions.HelpFunc = goptions.NewTemplatedHelpFunc(falsy.FALSE_HELP)
		flagset.Verbs["false"].HelpFunc = falseHelp
		flagset.Verbs["false"].PrintHelp(os.Stderr)
		os.Exit(1)
	case "true":
		if err == nil {
			runTrue(&options)
		}
		var trueHelp goptions.HelpFunc = goptions.NewTemplatedHelpFunc(truthy.TRUE_HELP)
		flagset.Verbs["true"].HelpFunc = trueHelp
		flagset.Verbs["true"].PrintHelp(os.Stderr)
		os.Exit(1)
	case "wc":
		// Yup: I really hate this.
		var options *pflag.FlagSet = pflag.NewFlagSet("wc", pflag.ContinueOnError)
		options.BoolP(wc.LineFlag, "l", false, "Count newlines")
		options.BoolP(wc.ByteFlag, "c", false, "Count bytes")
		options.BoolP(wc.WordFlag, "w", false, "Count words")
		options.BoolP(wc.CharFlag, "m", false, "Count characters")
		options.BoolP(wc.LongFlag, "L", false, "Print longest line length")
		if err := options.Parse(os.Args[2:]); err != nil {
			fmt.Fprint(os.Stderr, err)
			options.Usage()
		}
		if err := wc.WcMain(options); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		flagset.PrintHelp(os.Stderr)
		os.Exit(1)
	}
}
