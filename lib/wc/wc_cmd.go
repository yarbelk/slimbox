package wc

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"gitlab.com/yarbelk/slimbox/lib"
)

func init() {
	lib.RegisterFunction("wc")
}

type errors []error

func (es errors) Error() string {
	builder := strings.Builder{}
	for _, e := range es {
		builder.WriteString(e.Error())
		builder.WriteRune('\n')
	}
	return builder.String()
}

// ReadFile from a chan and stream out results.
// It was a closure over what the arguments are, but I want to pull it out to test
func ReadFile(options Options, fnChan <-chan string, resChan chan<- Results, errChan chan<- error, wg *sync.WaitGroup) {
	for filename := range fnChan {
		_, in, err := lib.ParseFiles(filename)
		if err != nil {
			resChan <- Results{Filename: filename}
			errChan <- err
		}
		if filename != "-" {
			defer in.Close()
		}

		fmt.Println(filename)
		results, err := WordCount(options, in)
		results.Filename = filename

		if err != nil {
			errChan <- err
		}
		resChan <- results
	}
	wg.Done()
}

// Main is the kickoff for the wc program.  so it can be compiled stand alone or as a subcommand
// runs as many workers as CPUs.  this probably should be tunable at compile time
func Main(options Options) error {
	// First section is setting up filenames; makeing sure we know if
	// We will read from stdin (if so, we just use one worker so its all sequential, as stdin can be
	// read from multiple times; so order matters)
	hasStdin := false
	if len(options.Files) == 0 {
		options.Files = []string{"-"}
		hasStdin = true
	} else {
		for _, fn := range options.Files {
			if fn == "-" {
				hasStdin = true
				break
			}
		}
	}

	// assemble the datastructures for collecting results
	order := make(map[string]int)
	r := make([]Results, len(options.Files))
	for i, filename := range options.Files {
		order[filename] = i
	}
	workers := runtime.NumCPU()
	if workers > len(options.Files) {
		workers = len(options.Files)
	}
	if hasStdin {
		workers = 1
	}

	fnChan := make(chan string)
	resChan := make(chan Results)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	errs := make(errors, 0)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go ReadFile(options, fnChan, resChan, errChan, &wg)
	}

	go func() {
		for _, filename := range options.Files {
			fnChan <- filename
		}
		close(fnChan)
	}()
	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()

	var max uint
readLoop:
	for {
		select {
		case c, ok := <-resChan:
			if !ok {
				break readLoop
			}
			r[order[c.Filename]] = c
			if c.Longest > max {
				max = c.GetMax(options)
			}
		case err, ok := <-errChan:
			if !ok {
				break readLoop
			}
			errs = append(errs, err)
		}

	}

	var err error
	if len(errs) > 0 {
		err = errs
	}
	resultsSet := ResultsSet{MaxNumber: max, Results: r}
	fmt.Print(resultsSet.Printf(options))
	return err
}

func BindFlagSet(wo *Options) *pflag.FlagSet {
	var wcFS *pflag.FlagSet = pflag.NewFlagSet("wc", pflag.ContinueOnError)
	wcFS.BoolVarP(&wo.Newlines, LineFlag, "l", false, "Count newlines")
	wcFS.BoolVarP(&wo.Bytes, ByteFlag, "c", false, "Count bytes")
	wcFS.BoolVarP(&wo.Words, WordFlag, "w", false, "Count words")
	wcFS.BoolVarP(&wo.Characters, CharFlag, "m", false, "Count characters")
	wcFS.BoolVarP(&wo.Longest, LongFlag, "L", false, "Print longest line length")
	setUsage(wcFS)
	return wcFS
}

func setUsage(wcFS *pflag.FlagSet) {
	wcFS.Usage = func() {

		fmt.Fprint(wcFS.Output(), `Usage: wc [OPTION]... [FILE]...
Print newline, word, and byte counts for each FILE, and a total line if
more than one FILE is specified.  A word is a non-zero-length sequence of
characters delimited by white space.

With no FILE, or when FILE is -, read standard input.  - can appear multiple
times in the list, and standard input will be read for each.

If there is no standard input, wc will run as many workers as there are CPUs
in your system and parallelize the work.

The options below may be used to select which counts are printed, always in
the following order: newline, word, character, byte, maximum line length.
`)
		wcFS.PrintDefaults()
	}
}
