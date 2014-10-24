package cat

import (
	"bytes"
	"strings"
	"testing"
)

var c *catOptions = &catOptions{}

// Catting an empty file should result in an empty response
func TestNoLineCat(t *testing.T) {
	var inputValue string = ""
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = ""

	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}

}

// A File with one line and a new line should be the same
func TestOneLineCat(t *testing.T) {
	var inputValue string = "Hello World!\n"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello World!\n"

	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}

// A file without a newline at the end should also have a newline
func TestOneLineCatNoNewline(t *testing.T) {
	var inputValue string = "Hello World!"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello World!\n"

	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}

// A with multiple lines should be output on multiiple lines
func TestMultiLineCat(t *testing.T) {
	var inputValue string = "Hello World!\n\nHello Gophers!\n"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello World!\n\nHello Gophers!\n"

	c.Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}

// Calling cat twice should concatinate the results
func TestMultiInputCat(t *testing.T) {
	var inputValueOne string = "Hello World!\n"
	var inputValueTwo string = "Hello Gophers!\n"
	var inputReaderOne = strings.NewReader(inputValueOne)
	var inputReaderTwo = strings.NewReader(inputValueTwo)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello World!\nHello Gophers!\n"

	c.Cat(inputReaderOne, &outputStream)
	c.Cat(inputReaderTwo, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}

// Calling cat twice should concatinate the results with a newline if it is missing
func TestMultiInputCatMissingNewline(t *testing.T) {
	var inputValueOne string = "Hello World!"
	var inputValueTwo string = "Hello Gophers!\n"
	var inputReaderOne = strings.NewReader(inputValueOne)
	var inputReaderTwo = strings.NewReader(inputValueTwo)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello World!\nHello Gophers!\n"

	c.Cat(inputReaderOne, &outputStream)
	c.Cat(inputReaderTwo, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}
}
