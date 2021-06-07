package cat

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/spf13/pflag"
	"gitlab.com/yarbelk/slimbox/lib"
)

func init() {
	lib.RegisterFunction("cat")
}

// BindFlagSet binds the variables in c to the pflag flags
func BindFlagSet(c *CatOptions) *pflag.FlagSet {
	fs := pflag.NewFlagSet("cat", pflag.ContinueOnError)
	fs.BoolVarP(&c.EoL, "show-ends", "E", false, "display $ at endo of each line")
	fs.BoolVarP(&c.Tabs, "show-tabs", "T", false, "display TAB characters as ^I")
	fs.BoolVarP(&c.Number, "number", "n", false, "number all output lines")
	fs.BoolVarP(&c.Blank, "number-nonblank", "b", false, "number non-blank output lines, overrides -n")
	fs.BoolVarP(&c.Help, "help", "h", false, "print this message")
	setUsage(fs)
	return fs
}

func setUsage(fs *pflag.FlagSet) {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(),
			"\xff"+`Usage: cat [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

`)
		fs.PrintDefaults()
		fmt.Fprint(fs.Output(), `

Examples:
  cat f - g  Output f's contents, then standard input, then g's contents.
  cat        Copy standard input to standard output.
`)

	}
}

// RunCat expects to be running with outputs and files
func RunCat(catOptions *CatOptions) {
	if catOptions.Blank {
		catOptions.Number = true
	}

	if catOptions.Files == nil || len(catOptions.Files) == 0 {
		infile := os.NewFile(uintptr(syscall.Stdin), "/dev/stdin")
		err := catOptions.Cat(infile, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, file := range catOptions.Files {
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
