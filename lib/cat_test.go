package cat

import (
	"bytes"
	"strings"
	"testing"
)

func TestOneLineCat(t *testing.T) {
	var inputValue string = "Hello World!\n"
	var inputReader = strings.NewReader(inputValue)
	var outputStream bytes.Buffer
	var expectedValue string = "Hello World!\n"

	Cat(inputReader, &outputStream)

	recievedValue := outputStream.String()

	if expectedValue != recievedValue {
		t.Fatalf("Expected %q, was %q", expectedValue, recievedValue)
	}

}
