package truthy

import (
	"os"
	"github.com/voxelbrain/goptions"
)

const (
	TRUE_HELP = "\xffDo nothing and succeed at it." +
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

type TrueOptions struct {
	Version   bool          `goptions:"--version, description='output version information and exit'"`
	Help      goptions.Help `goptions:"-h, --help, description='print this message'"`
}

func True() {
	os.Exit(0)
}
