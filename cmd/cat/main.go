package main

import (
	"os"

	"gitlab.com/yarbelk/slimbox/lib/cat"
)

func main() {
	catOptions := cat.NewCatOptions()
	catFS := cat.BindFlagSet(catOptions)
	err := catFS.Parse(os.Args[1:])
	if catOptions.Help {
		catFS.Usage()
		os.Exit(0)
	}

	if err == nil {
		catOptions.Files = catFS.Args()
		cat.RunCat(catOptions)
		os.Exit(0)
	}
}
