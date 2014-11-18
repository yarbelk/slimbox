package wc

import (
	"strings"
	"testing"
)

var wc *WCOptions = &WCOptions{}

// wc an empty file should result in a zero response
func TestEmptyFileWC(t *testing.T) {
	var inputValue string = ""
	var inputReader = strings.NewReader(inputValue)
	var filename string = "blank_file"

	fileStats := &FileStats{Filename: filename}
	wc.oneFileWC(inputReader, fileStats)

	if fileStats.words != 0 || fileStats.bytes != 0 || fileStats.lines != 0 {
		t.Fatalf("Expected zero in everything, got %v", fileStats)
	}

}

// wc a file with one word result in a one word count response
func TestOneWordWC(t *testing.T) {
	var inputValue string = "hello"
	var inputReader = strings.NewReader(inputValue)
	var filename string = "blank_file"

	fileStats := &FileStats{Filename: filename}
	wc.oneFileWC(inputReader, fileStats)

	if fileStats.words != 1 {
		t.Fatalf("expected one word, got %d", fileStats.words)
	}
	if fileStats.bytes != 5 {
		t.Fatalf("expected 5 bytes , got %d", fileStats.bytes)
	}
	if fileStats.lines != 0 {
		t.Fatalf("expected no lines , got %d", fileStats.lines)
	}

}
