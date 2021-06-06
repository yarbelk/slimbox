package wc_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"gitlab.com/yarbelk/slimbox/lib/wc"
)

func TestReadFileFailures(t *testing.T) {
	fnChan := make(chan string)
	resChan := make(chan wc.Results)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go wc.ReadFile(wc.DefaultOptions{Filenames: []string{"."}}, fnChan, resChan, errChan, &wg)
	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()

	fnChan <- "."
	close(fnChan)
	err := <-errChan
	res := <-resChan

	expected := wc.Results{Filename: "."}

	if err.Error() != "read .: is a directory" {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(expected, res) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, res)
		t.FailNow()
	}
}
