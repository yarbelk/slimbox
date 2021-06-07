package falsy

import (
	"os"

	"gitlab.com/yarbelk/slimbox/lib"
)

func init() {
	lib.RegisterFunction("false")
}

// False is only False.  it only exits badly
func False() {
	os.Exit(1)
}
