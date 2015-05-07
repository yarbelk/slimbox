package falsy

import (
	"os"
	"github.com/voxelbrain/goptions"
)

const (
	FALSE_HELP = "\xffDo nothing and fail at it." +
		"Usage: {{.Name}} [ignored command line arguments]" +
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
		"{{end}}\n"
)

type FalseOptions struct {
	Version   bool          `goptions:"--version, description='output version information and exit'"`
	Help      goptions.Help `goptions:"-h, --help, description='print this message'"`
}

func False() {
	os.Exit(1)
}
