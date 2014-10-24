package cat

import (
	"bytes"
	"strings"
	"testing"
)

// You should be able to add an '$' at the end of the line
func TestAddEOLChar(t *testing.T) {
	var inputValue string = "Hello\nWorld!\n"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello$\nWorld!$\n"

	var c *catOptions

	c = &catOptions{eol: true}

	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}

// Escape tab helper adds an escape character
func TestAddEscapedTabs(t *testing.T) {
	var inputValue string = "Hello\tWorld!\n"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello\\tWorld!\n"

	c = &catOptions{tabs: true}
	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}
