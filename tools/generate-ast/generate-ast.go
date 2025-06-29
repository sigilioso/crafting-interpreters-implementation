package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	types := []string{
		"Binary   : Left Expr[T], Operator tokens.Token, Right Expr[T]",
		"Grouping : Expression Expr[T]",
		"Literal  : Value any",
		"Unary    : Operator tokens.Token, Right Expr[T]",
	}

	defineAst("../../glox/expr", "Expr", types)

	// run go fmt
	cmd := exec.Command("go", "fmt")
	cmd.Dir = "../../glox/expr"
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func defineAst(outputDir string, basename string, types []string) {
	filePath := path.Join(outputDir, strings.ToLower(basename)+".go")
	_ = os.Remove(filePath)
	var code strings.Builder
	fmt.Fprintln(&code, "// Generated via tools/generate-ast")
	fmt.Fprintln(&code, "package expr")
	fmt.Fprintln(&code)
	fmt.Fprintf(&code, "import %q\n", "glox/tokens")
	fmt.Fprintln(&code)
	fmt.Fprintf(&code, "type %s[T any] interface {\n", basename)
	fmt.Fprintln(&code, "Accept(Visitor[T]) (T, error)")
	fmt.Fprintln(&code, "}")
	fmt.Fprintln(&code)
	for _, typeDef := range types {
		defineType(&code, typeDef)
	}
	defineVisitor(&code, types)

	if err := os.WriteFile(filePath, []byte(code.String()), 0644); err != nil {
		panic(err)
	}
}

func defineType(code *strings.Builder, typeDef string) {
	chunks := strings.Split(typeDef, ":")
	name := strings.TrimSpace(chunks[0])
	fields := strings.TrimSpace(chunks[1])

	fmt.Fprintf(code, "type %s[T any] struct {\n", name)
	for _, field := range strings.Split(fields, ", ") {
		fmt.Fprintln(code, field)
	}
	fmt.Fprintln(code, "}")
	fmt.Fprintln(code)
	fmt.Fprintf(code, "func (e %s[T]) Accept(v Visitor[T]) (T, error) {\n", name)
	fmt.Fprintf(code, "return v.VisitFor%s(e)\n", name)
	fmt.Fprintln(code, "}")

	fmt.Fprintln(code)

}

func defineVisitor(code *strings.Builder, types []string) {

	fmt.Fprintln(code, "type Visitor[T any] interface {")
	for _, typeDef := range types {
		chunks := strings.Split(typeDef, ":")
		typeName := strings.TrimSpace(chunks[0])
		fmt.Fprintf(code, "VisitFor%s(%s[T]) (T, error)\n", typeName, typeName)
	}
	fmt.Fprintln(code, "}")
	fmt.Fprintln(code)
}
