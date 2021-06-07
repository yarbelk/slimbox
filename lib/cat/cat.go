package cat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type CatOptions struct {
	EoL       bool
	Tabs      bool
	Number    bool
	Blank     bool
	Help      bool
	Files     []string
	curNumber int
}

type lineMutator func(io.Reader) ([]byte, error)

func escapeTabs(inputStream io.Reader) ([]byte, error) {
	bufferedReader := bufio.NewScanner(inputStream)
	outputStream := new(bytes.Buffer)
	bufferedReader.Split(bufio.ScanBytes)
	for bufferedReader.Scan() {
		aByte := bufferedReader.Bytes()
		if aByte[0] != byte('\t') {
			outputStream.Write(aByte)
		} else {
			outputStream.Write([]byte{'^', 'I'})
		}
	}
	if err := bufferedReader.Err(); err != nil {
		return nil, err
	}
	return outputStream.Bytes(), nil
}

func (c *CatOptions) prependLineNumber(inputStream io.Reader) ([]byte, error) {
	bufferedReader := bufio.NewScanner(inputStream)
	outputStream := new(bytes.Buffer)
	if !c.Blank {
		numbers := fmt.Sprintf("%6d  ", c.curNumber)
		c.curNumber++
		outputStream.Write([]byte(numbers))
	}
	for bufferedReader.Scan() {
		if c.Blank {
			numbers := fmt.Sprintf("%6d  ", c.curNumber)
			c.curNumber++
			outputStream.Write([]byte(numbers))
		}
		outputStream.Write(bufferedReader.Bytes())
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

// NewCatOptions returns a pointer to a default CatOptions struct, which means curNumber is 1
func NewCatOptions() *CatOptions {
	c := &CatOptions{}
	c.curNumber = 1
	return c

}

func (c *CatOptions) Cat(originalInputStream io.Reader, outputStream io.Writer) error {
	var (
		line []byte
		err  error
	)

	bufferedReader := bufio.NewScanner(originalInputStream)
	for bufferedReader.Scan() {
		line = bufferedReader.Bytes()
		if c.EoL {
			line, err = appendEOL(bytes.NewBuffer(line))
			if err != nil {
				break
			}
		}
		if c.Tabs {
			line, err = escapeTabs(bytes.NewBuffer(line))
			if err != nil {
				break
			}
		}
		if c.Number {
			line, err = c.prependLineNumber(bytes.NewBuffer(line))
			if err != nil {
				break
			}
		}

		outputStream.Write(line)
		outputStream.Write([]byte{byte('\n')})
	}

	// Things that act on the final line (even if it is added by cat) as well
	if c.Number {
		line, err = c.prependLineNumber(bytes.NewBuffer([]byte("")))
		outputStream.Write(line)
		if err != nil {
			return err
		}
	}
	if err := bufferedReader.Err(); err != nil {
		return err
	}
	return nil
}
