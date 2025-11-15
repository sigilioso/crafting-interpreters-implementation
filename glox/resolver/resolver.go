package resolver

import (
	"glox/expr"
	"glox/interpreter"
	"glox/stmt"
	"glox/tokens"
)

type Resolver struct {
	Interpreter interpreter.Interpreter
	scopes      Stack[map[string]bool]
}

// VisitForBlock implements stmt.Visitor.
func (r *Resolver) VisitForBlock(s stmt.Block[any]) (any, error) {
	r.beginScope()
	if err := r.ResolveStatements(s.Statements); err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil

}

// VisitForExpression implements stmt.Visitor.
func (r *Resolver) VisitForExpression(stmt.Expression[any]) (any, error) {
	panic("unimplemented")
}

// VisitForFunction implements stmt.Visitor.
func (r *Resolver) VisitForFunction(stmt.Function[any]) (any, error) {
	panic("unimplemented")
}

// VisitForIf implements stmt.Visitor.
func (r *Resolver) VisitForIf(stmt.If[any]) (any, error) {
	panic("unimplemented")
}

// VisitForPrint implements stmt.Visitor.
func (r *Resolver) VisitForPrint(stmt.Print[any]) (any, error) {
	panic("unimplemented")
}

// VisitForReturn implements stmt.Visitor.
func (r *Resolver) VisitForReturn(stmt.Return[any]) (any, error) {
	panic("unimplemented")
}

// VisitForVar implements stmt.Visitor.
func (r *Resolver) VisitForVar(s stmt.Var[any]) (any, error) {
	r.declare(s.Name)
	if s.Initializer != nil {
		panic("TODO")
		// r.ResolveStatement(s.Initializer)
	}
	r.define(s.Name)
	return nil, nil
}

// VisitForWhile implements stmt.Visitor.
func (r *Resolver) VisitForWhile(stmt.While[any]) (any, error) {
	panic("unimplemented")
}

func (r *Resolver) VisitForVariable(v expr.Variable[any]) (any, error) {
	panic("unimplemented")
}

func (r *Resolver) ResolveStatements(statements []stmt.Stmt[any]) error {
	for _, statement := range statements {
		if err := r.ResolveStatement(statement); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) ResolveStatement(statement stmt.Stmt[any]) error {
	_, err := statement.Accept(r)
	return err
}

func (r *Resolver) beginScope() {
	r.scopes.Push(map[string]bool{})
}

func (r *Resolver) endScope() {
	_ = r.scopes.Pop()
}

func (r *Resolver) declare(name tokens.Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek()
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name tokens.Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expression expr.Expr[any], name tokens.Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope := r.scopes.Get(i)
		if _, containsKey := scope[name.Lexeme]; containsKey {
			panic("TODO")
			//r.Interpreter.Resolve(expression, r.scopes.Size()-1 -i)
			// return
		}
	}
}

// Stack is a simple stack implementation
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (item T) {
	if l := len(s.items); l > 0 {
		item = s.items[l-1]
		s.items = s.items[:l-1]
		return item
	}
	return item
}

func (s *Stack[T]) Peek() (item T) {
	if l := len(s.items); l > 0 {
		item = s.items[l-1]
		return item
	}
	return item
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) Get(i int) T {
	return s.items[i]
}
