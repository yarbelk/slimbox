package lib

import "strings"

var registry Functions = make(Functions, 0)

type Functions []string

func (f Functions) String() string {
	builder := strings.Builder{}
	for i, fn := range f {
		if i != 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fn)
	}
	return builder.String()
}

func RegisterFunction(function string) {
	for _, fn := range registry {
		if fn == function {
			return
		}
	}
	registry = append(registry, function)
}

// RegisteredFunctions returns a copy of the registered functions
func RegisteredFunctions() Functions {
	return registry[:]
}
