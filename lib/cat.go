package cat

import (
	"bufio"
	"io"
)

func Cat(inputStream io.Reader, outputStream io.Writer) error {
	bufferedReader := bufio.NewScanner(inputStream)

	for bufferedReader.Scan() {
		outputStream.Write(bufferedReader.Bytes())
		outputStream.Write([]byte{byte('\n')})
	}
	if err := bufferedReader.Err(); err != nil {
		return err
	}

	return nil

}
