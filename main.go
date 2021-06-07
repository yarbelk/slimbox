package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/yarbelk/slimbox/lib/falsy"
	"github.com/yarbelk/slimbox/lib/truthy"
	"gitlab.com/yarbelk/slimbox/lib"
	"gitlab.com/yarbelk/slimbox/lib/cat"
	"gitlab.com/yarbelk/slimbox/lib/wc"
)

var catOptions cat.CatOptions

var numberLines, numberNonBlankLines, help bool
var lineNumber int = 1

func version() {
	os.Stderr.WriteString("Nice try - no versions yet\n")
	os.Exit(0)
}
func init() {
	// Complier seems to mess this up...
	lib.RegisterFunction("true")
	lib.RegisterFunction("false")
}

var _ func() = truthy.True
var _ func() = falsy.False

func main() {

	// This is all placeholder until i migrate true, false and cat to pflag (fix argument ordering)
	var verb string
	if len(os.Args) > 1 {
		verb = os.Args[1]
	}
	switch verb {
	case "cat":
		catOptions := cat.NewCatOptions()
		catFS := cat.BindFlagSet(catOptions)
		err := catFS.Parse(os.Args[2:])
		if catOptions.Help {
			catFS.Usage()
			os.Exit(0)
		}

		if err == nil {
			catOptions.Files = catFS.Args()
			cat.RunCat(catOptions)
			os.Exit(0)
		}

		os.Exit(1)
	case "false":
		falsy.False()
	case "true":
		truthy.True()
	case "wc":
		wcOptions := wc.Options{}
		wcFS := wc.BindFlagSet(&wcOptions)
		if err := wcFS.Parse(os.Args[2:]); err != nil {
			if err != pflag.ErrHelp {
				fmt.Fprint(os.Stderr, err)
				wcFS.Usage()
			}
			os.Exit(1)
		}
		wcOptions.Files = wcFS.Args()[:]
		if err := wc.Main(wcOptions); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, `
Usage: slimbox [function [arguments]...]
   or: function [arguments]...

slimbox is %s WIP.  Its also probably slower than any other compiled implementation.
I mean; take alook at how I'm making that really bold; its silly.
The currently supported functions are:

%s
`, "\033[1mreally\033[0m", lib.RegisteredFunctions())
		os.Exit(1)
	}
}
