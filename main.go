package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/voxelbrain/goptions"
	"github.com/yarbelk/slimbox/lib"
	"github.com/yarbelk/slimbox/lib/cat"
	"github.com/yarbelk/slimbox/lib/falsy"
	"github.com/yarbelk/slimbox/lib/truthy"
)

var catOptions cat.CatOptions

var numberLines, numberNonBlankLines, help bool
var lineNumber int = 1

type Options struct {
	Verb  goptions.Verbs
	Cat   cat.CatOptions     `goptions:"cat"`
	False falsy.FalseOptions `goptions:"false"`
	True  truthy.TrueOptions `goptions:"true"`
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

	if err != nil {
		if options.Verb == "cat" {
			var catHelp goptions.HelpFunc = goptions.NewTemplatedHelpFunc(cat.CAT_HELP)
			flagset.Verbs["cat"].HelpFunc = catHelp
			flagset.Verbs["cat"].PrintHelp(os.Stderr)
			os.Exit(1)
		} else if options.Verb == "false" {
			var falseHelp goptions.HelpFunc = goptions.NewTemplatedHelpFunc(falsy.FALSE_HELP)
			flagset.Verbs["false"].HelpFunc = falseHelp
			flagset.Verbs["false"].PrintHelp(os.Stderr)
			os.Exit(1)
		} else if options.Verb == "True" {
			var trueHelp goptions.HelpFunc = goptions.NewTemplatedHelpFunc(truthy.TRUE_HELP)
			flagset.Verbs["true"].HelpFunc = trueHelp
			flagset.Verbs["true"].PrintHelp(os.Stderr)
			os.Exit(1)
		}

		flagset.PrintHelp(os.Stderr)
		os.Exit(1)
	}

	if options.Verb == "cat" {
		runCat(&options)
	} else if options.Verb == "false" {
		runFalse(&options)
	} else if options.Verb == "true" {
		runTrue(&options)
	} else {
		flagset.PrintHelp(os.Stderr)
	}
}
