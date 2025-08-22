package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	types_expr := []string{
		"Assign	  : Name tokens.Token, Value Expr[T]",
		"Binary   : Left Expr[T], Operator tokens.Token, Right Expr[T]",
		"Grouping : Expression Expr[T]",
		"Literal  : Value any",
		"Unary    : Operator tokens.Token, Right Expr[T]",
		"Variable : Name tokens.Token",
	}
	defineAst("../../glox/expr", "Expr", types_expr)

	types_stmt := []string{
		"Block		: Statements []Stmt[T]",
		"Expression	: Expression expr.Expr[T]",
		"Print		: Expression expr.Expr[T]",
		"Var		: Name tokens.Token, Initializer expr.Expr[T]",
	}
	defineAst("../../glox/stmt", "Stmt", types_stmt)

	// format go code
	dirs := []string{
		"../../glox/expr",
		"../../glox/stmt",
	}
	for _, dir := range dirs {
		cmd := exec.Command("goimports", "-w", ".")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}

func defineAst(outputDir string, basename string, types []string) {
	packageName := path.Base(outputDir)
	filePath := path.Join(outputDir, strings.ToLower(basename)+".go")
	_ = os.Remove(filePath)
	var code strings.Builder
	fmt.Fprintln(&code, "// Generated via tools/generate-ast")
	fmt.Fprintf(&code, "package %s\n", packageName)
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
