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

	var c *CatOptions

	c = &CatOptions{EoL: true}

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
	var expectedValue string = "Hello^IWorld!\n"

	c = &CatOptions{Tabs: true}
	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}

// number the lines
func TestAddLineNumbers(t *testing.T) {
	var inputValue string = "Hello\nWorld!\n"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "     1  Hello\n     2  World!\n     3  "

	c = NewCatOptions()
	c.Number = true
	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}
