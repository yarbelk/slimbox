package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/spf13/pflag"
	"gitlab.com/yarbelk/slimbox/lib/wc"
)

var (
	options      *pflag.FlagSet = pflag.NewFlagSet("wc", pflag.ContinueOnError)
	lines        *bool          = options.BoolP("lines", "l", false, "Count newlines")
	bytes        *bool          = options.BoolP("bytes", "c", false, "Count bytes")
	words        *bool          = options.BoolP("words", "w", false, "Count words")
	characters   *bool          = options.BoolP("chars", "m", false, "Count characters")
	printLongist *bool          = options.BoolP("max-line-length", "L", false, "Print longest line length")
)

func WcMain() error {
	var resultsSet wc.ResultsSet
	if err := options.Parse(os.Args[1:]); err != nil {
		fmt.Fprint(os.Stderr, err)
		options.Usage()
	}
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

	order := make(map[string]int)
	r := make([]wc.Results, len(filenames))
	for i, filename := range filenames {
		order[filename] = i
	}
	workers := runtime.NumCPU()
	if hasStdin {
		workers = 1
	}
	fnChan := make(chan string)
	resChan := make(chan wc.Results)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for filename := range fnChan {
				fmt.Println(filename)
				var in *os.File
				var err error
				if filename == "-" {
					in = os.Stdin
				} else {
					in, err = os.Open(filename)
					if err != nil {
						errChan <- err
					}
					defer in.Close()
				}
				results, err := wc.WordCount(options, in)
				if err != nil {
					errChan <- err
				}
				if filename != "-" {
					results.Filename = filename
				}
				resChan <- results

			}
			wg.Done()
		}()
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
	for c := range resChan {
		r[order[c.Filename]] = c
		if c.Longest > max {
			max = c.GetMax(options)
		}
	}
	resultsSet = wc.ResultsSet{MaxNumber: max, Results: r}
	fmt.Print(resultsSet.Printf(options))
	return nil

}

func main() {
	WcMain()
}
