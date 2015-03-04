package falsy

import (
	"os"
	"github.com/voxelbrain/goptions"
)

const (
	FALSE_HELP = "\xffExit with a status code indicating failure." +
		"Usage: {{.Name}} [ignored command line arguments]" +
		"\n" +
		"or: {{.Name}} OPTION" +
		"\n" +
		"{{end}}" +
		"{{end}}" +
		"{{end}}" +
		"{{end}}"
)

type FalseOptions struct {
	Version   bool          `goptions:"--version, description='output version information and exit'"`
	Help      goptions.Help `goptions:"-h, --help, description='print this message'"`
}

func False() {
	os.Exit(1)
}
