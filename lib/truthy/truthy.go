package truthy

import (
	"os"

	"gitlab.com/yarbelk/slimbox/lib"
)

func init() {
	lib.RegisterFunction("true")
}

// True is only true.  it only exits well
func True() {
	os.Exit(0)
}
