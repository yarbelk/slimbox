package cat

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/voxelbrain/goptions"
	"io"
)

const (
	CAT_HELP = "\xffConcatenate file[s] or standard input to standard output\n" +
		"Usage: {{.Name}} [global options] {{with .Verbs}}<verb> [verb options]{{end}}\n" +
		"\n" +
		"Global options:" +
		"{{range .Flags}}" +
		"\n\t" +
		"\t{{with .Short}}" + "-{{.}}," + "{{end}}" +
		"\t{{with .Long}}" + "--{{.}}" + "{{end}}" +
		"\t{{.Description}}" +
		"{{with .DefaultValue}}" +
		" (default: {{.}})" +
		"{{end}}" +
		"{{if .Obligatory}}" +
		" (*)" +
		"{{end}}" +
		"{{end}}" +
		"\xff\n\n{{with .Verbs}}Verbs:\xff" +
		"{{range .}}" +
		"\xff\n    {{.Name}}:\xff" +
		"{{range .Flags}}" +
		"\n\t" +
		"\t{{with .Short}}" + "-{{.}}," + "{{end}}" +
		"\t{{with .Long}}" + "--{{.}}" + "{{end}}" +
		"\t{{.Description}}" +
		"{{with .DefaultValue}}" +
		" (default: {{.}})" +
		"{{end}}" +
		"{{if .Obligatory}}" +
		" (*)" +
		"{{end}}" +
		"{{end}}" +
		"{{end}}" +
		"{{end}}" +
		"\n" +
		"With no FILE, or when FILE is -, read standard input.\n\n" +
		"Examples:\n" +
		"  cat f - g  Output f's contents, then standard input, then g's contents.\n" +
		"  cat        Copy standard input to standard output.\n"
)

type CatOptions struct {
	EoL       bool          `goptions:"-E, --show-ends, description='display $ at end of each line'"`
	Tabs      bool          `goptions:"-T, --show-tabs, description='display TAB characters as ^I'"`
	Number    bool          `goptions:"-n, --number, description='number all output lines'"`
	Blank     bool          `goptions:"-b, --number-nonblank, description='number non-blank output lines, overrides -n'"`
	Help      goptions.Help `goptions:"-h, --help, description='print this message'"`
	Files     goptions.Remainder
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
