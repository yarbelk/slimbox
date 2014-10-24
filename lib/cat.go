package cat

import (
	"bufio"
	"io"
)

type catOptions struct {
	eol bool
}

func (c *catOptions) Cat(inputStream io.Reader, outputStream io.Writer) error {
	bufferedReader := bufio.NewScanner(inputStream)

	for bufferedReader.Scan() {
		outputStream.Write(bufferedReader.Bytes())
		if c.eol {
			outputStream.Write([]byte{byte('$')})
		}
		outputStream.Write([]byte{byte('\n')})
	}
	if err := bufferedReader.Err(); err != nil {
		return err
	}
	return nil
}
