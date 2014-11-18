package lib

import (
	"os"
	"syscall"
)

func ParseFiles(filename string) (string, *os.File, error) {
	if filename == "-" {
		return "/dev/stdin", os.NewFile(uintptr(syscall.Stdin), "/dev/stdin"), nil
	}
	fd, err := os.Open(filename)
	return filename, fd, err

}
