package errors

import (
	"fmt"
	"glox/tokens"
	"os"
)

var errorFound = false
var runtimeErrorFound = false

func AtLine(line int, message string) {
	Report(line, "", message)
}

func AtToken(token tokens.Token, message string) {
	if token.TokenType == tokens.Eof {
		Report(token.Line, " at end", message)
	} else {
		Report(token.Line, fmt.Sprintf(" at '%s'", token.Lexeme), message)
	}
}

func Report(line int, where string, message string) {
	errorFound = true
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s", line, where, message)
}

func ReportRuntimeError(token tokens.Token, message string) {
	runtimeErrorFound = true
	fmt.Fprintf(os.Stderr, "%s\n[line %d]", message, token.Line)
}

func ErrorFound() bool {
	return errorFound
}

func RuntimeErrorFound() bool {
	return runtimeErrorFound
}

func ResetError() {
	errorFound = false
}
