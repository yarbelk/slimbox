package cat

import (
	"bufio"
	"bytes"
	"io"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(4)
}

type catOptions struct {
	eol  bool
	tabs bool
}

func escapeTabs(inputStream io.Reader) ([]byte, error) {
	bufferedReader := bufio.NewScanner(inputStream)
	outputStream := new(bytes.Buffer)
	bufferedReader.Split(bufio.ScanBytes)
	for bufferedReader.Scan() {
		aByte := bufferedReader.Bytes()
		if aByte[0] != byte('\t') {
			outputStream.Write(aByte)
		} else {
			outputStream.Write([]byte{'\\', 't'})
		}
	}
	if err := bufferedReader.Err(); err != nil {
		return nil, err
	}
	return outputStream.Bytes(), nil
}

func appendEOL(inputStream io.Reader) ([]byte, error) {
	bufferedReader := bufio.NewScanner(inputStream)
	outputStream := new(bytes.Buffer)
	for bufferedReader.Scan() {
		outputStream.Write(bufferedReader.Bytes())
		outputStream.Write([]byte{byte('$')})
	}
	if err := bufferedReader.Err(); err != nil {
		return nil, err
	}
	return outputStream.Bytes(), nil
}

func (c *catOptions) Cat(originalInputStream io.Reader, outputStream io.Writer) error {
	var (
		line []byte
		err  error
	)
	var inputStream = originalInputStream

	bufferedReader := bufio.NewScanner(inputStream)
	for bufferedReader.Scan() {
		line = bufferedReader.Bytes()
		if c.eol {
			line, err = appendEOL(bytes.NewBuffer(line))
			if err != nil {
				break
			}
		}
		if c.tabs {
			line, err = escapeTabs(bytes.NewBuffer(line))
			if err != nil {
				break
			}
		}

		outputStream.Write(line)
		outputStream.Write([]byte{byte('\n')})
	}
	if err := bufferedReader.Err(); err != nil {
		return err
	}
	return nil
}
