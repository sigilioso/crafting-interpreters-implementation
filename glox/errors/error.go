package errors

import (
	"fmt"
	"os"
)

var errorFound = false

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	errorFound = true
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s", line, where, message)
}

func ErrorFound() bool {
	return errorFound
}

func ResetError() {
	errorFound = false
}
