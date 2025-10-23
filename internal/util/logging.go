package util

import (
	"fmt"
	"os"
)

func Debugf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[debug] "+format, args...)
}