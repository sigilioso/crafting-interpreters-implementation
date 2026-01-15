package resolver

import (
	"glox/errors"
	"glox/expr"
	"glox/interpreter"
	"glox/stmt"
	"glox/tokens"
)

type Resolver struct {
	Interpreter         *interpreter.Interpreter
	scopes              Stack[map[string]bool]
	currentFunctionType FunctionType
	currentClassType    ClassType
}

func NewResolver(i *interpreter.Interpreter) Resolver {
	return Resolver{Interpreter: i, scopes: Stack[map[string]bool]{}, currentFunctionType: FunctionTypeNone, currentClassType: ClassTypeNone}
}

func (r *Resolver) VisitForBlock(s *stmt.Block[any]) (any, error) {
	r.beginScope()
	if err := r.ResolveStatements(s.Statements); err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil

}

func (r *Resolver) VisitForExpression(e *stmt.Expression[any]) (any, error) {
	return nil, r.resolveExpr(e.Expression)
}

func (r *Resolver) VisitForFunction(f *stmt.Function[any]) (any, error) {
	r.declare(f.Name)
	r.define(f.Name)
	if err := r.resolveFunction(f, FunctionTypeFunction); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitForIf(s *stmt.If[any]) (any, error) {
	if err := r.resolveExpr(s.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(s.ThenBranch); err != nil {
		return nil, err
	}
	if s.ElseBranch != nil {
		if err := r.resolveStmt(s.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitForPrint(s *stmt.Print[any]) (any, error) {
	return nil, r.resolveExpr(s.Expression)
}

func (r *Resolver) VisitForReturn(s *stmt.Return[any]) (any, error) {
	if r.currentFunctionType == FunctionTypeNone {
		errors.AtToken(s.Keyword, "Can't return from top-level code.")
		return nil, nil
	}
	if s.Value != nil && r.currentFunctionType == FunctionTypeInitializer {
		errors.AtToken(s.Keyword, "Can't return a value from an initializer.")
	}
	if s.Value != nil {
		return nil, r.resolveExpr(s.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitForVar(s *stmt.Var[any]) (any, error) {
	r.declare(s.Name)
	if s.Initializer != nil {
		if err := r.resolveExpr(s.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(s.Name)
	return nil, nil
}

func (r *Resolver) VisitForWhile(s *stmt.While[any]) (any, error) {
	if err := r.resolveExpr(s.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(s.Body); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitForGet(g *expr.Get[any]) (any, error) {
	return nil, r.resolveExpr(g.Object)
}

func (r *Resolver) VisitForSet(s *expr.Set[any]) (any, error) {
	if err := r.resolveExpr(s.Value); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(s.Object); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitForClass(c *stmt.Class[any]) (any, error) {
	enclosingClassType := r.currentClassType
	r.currentClassType = ClassTypeClass
	r.declare(c.Name)
	r.define(c.Name)

	if c.SuperClass != nil {
		r.beginScope()
		r.scopes.Peek()["super"] = true
		if c.SuperClass.Name == c.Name {
			errors.AtToken(c.SuperClass.Name, "A class can't inherit from itself.")
		}
		r.currentClassType = ClassTypeSubclass
		if err := r.resolveExpr(c.SuperClass); err != nil {
			return nil, err
		}
	}

	r.beginScope()
	scope := r.scopes.Peek()
	scope["this"] = true

	for _, method := range c.Methods {
		declaration := FunctionTypeMethod
		if method.Name.Lexeme == "init" {
			declaration = FunctionTypeInitializer
		}
		if err := r.resolveFunction(method, declaration); err != nil {
			return nil, nil
		}
	}
	r.endScope()
	if c.SuperClass != nil {
		r.endScope()
	}
	r.currentClassType = enclosingClassType
	return nil, nil
}

func (r *Resolver) VisitForSuper(t *expr.Super[any]) (any, error) {
	if r.currentClassType == ClassTypeNone {
		errors.AtToken(t.Keyword, "Can't use 'super' outside of a class.")
		return nil, nil
	}
	if r.currentClassType != ClassTypeSubclass {
		errors.AtToken(t.Keyword, "Can't use 'super' in a class with no superclass.")
		return nil, nil
	}
	r.resolveLocal(t, t.Keyword)
	return nil, nil
}

func (r *Resolver) VisitForThis(t *expr.This[any]) (any, error) {
	if r.currentClassType == ClassTypeNone {
		errors.AtToken(t.Keyword, "Can't use 'this' outside of a class.")
		return nil, nil
	}
	r.resolveLocal(t, t.Keyword)
	return nil, nil
}

func (r *Resolver) VisitForVariable(v *expr.Variable[any]) (any, error) {
	if !r.scopes.IsEmpty() {
		if inScope, exits := r.scopes.Peek()[v.Name.Lexeme]; exits && !inScope {
			errors.AtToken(v.Name, "Can't read local variable in its own initializer.")
			return nil, nil
		}
	}
	r.resolveLocal(v, v.Name)
	return nil, nil
}

func (r *Resolver) ResolveStatements(statements []stmt.Stmt[any]) error {
	for _, statement := range statements {
		if err := r.ResolveStatement(statement); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) VisitForBinary(binary *expr.Binary[any]) (any, error) {
	if err := r.resolveExpr(binary.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(binary.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitForCall(c *expr.Call[any]) (any, error) {
	if err := r.resolveExpr(c.Callee); err != nil {
		return nil, err
	}
	for _, arg := range c.Arguments {
		if err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitForGrouping(grouping *expr.Grouping[any]) (any, error) {
	return nil, r.resolveExpr(grouping.Expression)
}

func (r *Resolver) VisitForLiteral(literal *expr.Literal[any]) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitForLogical(l *expr.Logical[any]) (any, error) {
	if err := r.resolveExpr(l.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(l.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitForAssign(a *expr.Assign[any]) (any, error) {
	if err := r.resolveExpr(a.Value); err != nil {
		return nil, err
	}
	r.resolveLocal(a, a.Name)
	return nil, nil
}

func (r *Resolver) VisitForUnary(unary *expr.Unary[any]) (any, error) {
	return nil, r.resolveExpr(unary.Right)
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
	if _, exits := scope[name.Lexeme]; exits {
		errors.AtToken(name, "Already a variable with this name in this scope.")
	}
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
			dept := r.scopes.Size() - 1 - i
			r.Interpreter.Resolve(expression, dept)
			return
		}
	}
}

func (r *Resolver) resolveExpr(v expr.Expr[any]) error {
	_, err := v.Accept(r)
	return err
}

func (r *Resolver) resolveStmt(s stmt.Stmt[any]) error {
	_, err := s.Accept(r)
	return err
}

func (r *Resolver) resolveStmtList(l []stmt.Stmt[any]) error {
	for _, s := range l {
		if err := r.resolveStmt(s); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(f *stmt.Function[any], functionType FunctionType) error {
	enclosingFunctionType := r.currentFunctionType
	r.currentFunctionType = functionType

	r.beginScope()
	for _, param := range f.Params {
		r.declare(param)
		r.define(param)
	}
	if err := r.resolveStmtList(f.Body); err != nil {
		return err
	}
	r.endScope()
	r.currentFunctionType = enclosingFunctionType
	return nil
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
