package wc

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

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
		var in *os.File
		var err error
		if filename == "-" {
			in = os.Stdin
		} else {
			in, err = os.Open(filename)
			if err != nil {
				resChan <- Results{Filename: filename}
				errChan <- err
			}
			defer in.Close()
		}
		results, err := WordCount(options, in)
		if err != nil {
			errChan <- err
		}
		if filename != "-" {
			results.Filename = filename
		}
		resChan <- results

	}
	wg.Done()

}

// WcMain is the kickoff for the wc program.  so it can be compiled stand alone or as a subcommand
// runs as many workers as CPUs.  this probably should be tunable at compile time
func WcMain(options Options) error {
	// First section is setting up filenames; makeing sure we know if
	// We will read from stdin (if so, we just use one worker so its all sequential, as stdin can be
	// read from multiple times; so order matters)
	filenames := options.Args()
	hasStdin := false
	if len(filenames) == 0 {
		filenames = []string{"-"}
		hasStdin = true
	} else {
		for _, fn := range filenames {
			if fn == "-" {
				hasStdin = true
				break
			}
		}
	}

	// assemble the datastructures for collecting results
	order := make(map[string]int)
	r := make([]Results, len(filenames))
	for i, filename := range filenames {
		order[filename] = i
	}
	workers := runtime.NumCPU()
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
		for _, filename := range filenames {
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
