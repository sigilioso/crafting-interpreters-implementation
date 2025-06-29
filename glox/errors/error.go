package errors

import (
	"fmt"
	"glox/tokens"
	"os"
)

var errorFound = false
var runtimeErrorFound = false

type RuntimeError struct {
	token   tokens.Token
	message string
}

func (e *RuntimeError) Error() string {
	return e.message
}

func (e *RuntimeError) Token() tokens.Token {
	return e.token
}

func NewRuntimeError(token tokens.Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

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

func ReportRuntimeError(e *RuntimeError) {
	runtimeErrorFound = true
	fmt.Fprintf(os.Stderr, "%s\n[line %d]", e.message, e.token.Line)
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
