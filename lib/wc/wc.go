package wc

import (
	"bufio"
	"fmt"
	"github.com/voxelbrain/goptions"
	"github.com/yarbelk/slimbox/lib"
	"io"
	"log"
)

const (
	WC_HELP = "Usage: {{.Name}} [global options] {{with .Verbs}}<verb> [verb options]{{end}}\n" +
		"\xffPrint newline, word, and byte counts for each FILE, and a total line if\n" +
		"more than one FILE is specified.  With no FILE, or when FILE is -,\n" +
		"read standard input.  A word is a non-zero-length sequence of characters\n" +
		"delimited by white space.\n" +
		"The options below may be used to select which counts are printed, always in\n" +
		"the following order: newline, word, character, byte, maximum line length.\n" +
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
		"{{end}}"
)

type FileStats struct {
	Filename string // filename, if any
	Offset   int    // byte offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count per line)
	bytes    int
	words    int
	chars    int
	lines    int
}

type WCOptions struct {
	Bytes         bool          `goptions:"-c, --bytes, description='print the byte counts'"`
	Chars         bool          `goptions:"-T, --show-tabs, description='print the character counts'"`
	Lines         bool          `goptions:"-n, --number, description='print the newline counts'"`
	Words         bool          `goptions:"-n, --number, description='print the word counts'"`
	MaxLineLength int           `goptions:"-b, --number-nonblank, description='print the length of the longest line'"`
	Help          goptions.Help `goptions:"-h, --help, description='print this message'"`
	Files         goptions.Remainder
}

func fmtLine(fs *FileStats) string {
	return fmt.Sprintf("%d %d %d %s", fs.lines, fs.words, fs.bytes, fs.Filename)
}

func (wc *WCOptions) WC(originalInputStream io.ByteReader, outputStream io.Writer) error {
	var fileStats []*FileStats = make([]*FileStats, len(wc.Files))
	for i, file := range wc.Files {
		func() {
			filename, fi, err := lib.ParseFiles(file) //TODO move to the root lib dir and rename
			if err != nil {
				log.Fatal(err)
			}
			defer fi.Close()
			fileStats[i] = &FileStats{Filename: filename}
			err = wc.oneFileWC(fi, fileStats[i])
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	return nil
}

func (wc *WCOptions) oneFileWC(originalInputStream io.Reader, filestat *FileStats) error {
	var line []byte

	bufferedReader := bufio.NewScanner(originalInputStream)
	for bufferedReader.Scan() {
		line = bufferedReader.Bytes()
		filestat.bytes += len(line) + 1
	}
	return nil
}
