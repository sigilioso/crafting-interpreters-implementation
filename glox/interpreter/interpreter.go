package interpreter

import (
	"fmt"
	"glox/environment"
	"glox/errors"
	"glox/expr"
	"glox/stmt"
	"glox/tokens"
	"strconv"
)

type Expr = expr.Expr[any]
type Binary = expr.Binary[any]
type Literal = expr.Literal[any]
type Unary = expr.Unary[any]
type Grouping = expr.Grouping[any]
type Variable = expr.Variable[any]
type Assign = expr.Assign[any]
type ExprVisitor = expr.Visitor[any]

type StmtExpression = stmt.Expression[any]
type StmtPrint = stmt.Print[any]
type StmtVisitor = stmt.Visitor[any]
type StmtVar = stmt.Var[any]
type StmtBlock = stmt.Block[any]
type Stmt = stmt.Stmt[any]

type Interpreter struct {
	env *environment.Environment
}

func New() Interpreter {
	return Interpreter{env: environment.New(nil)}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		if _, err := i.execute(statement); err != nil {
			e := err.(*errors.RuntimeError)
			errors.ReportRuntimeError(e)
			return
		}
	}
}

func (i *Interpreter) interpret(expression Expr) (string, error) {
	v, err := i.evaluate(expression)
	if err != nil {
		return "", err
	}
	return stringify(v), nil
}

func (i *Interpreter) execute(stmt Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, env *environment.Environment) error {
	previous := i.env
	i.env = env // Use the block's environment
	defer func() {
		i.env = previous // Switch back to previous environment when the block is over
	}()
	for _, statement := range statements {
		if _, err := i.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitForExpression(e StmtExpression) (any, error) {
	if _, err := i.evaluate(e.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitForPrint(p StmtPrint) (any, error) {
	v, err := i.evaluate(p.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(stringify(v))
	return nil, nil
}

func (i *Interpreter) VisitForBlock(b StmtBlock) (any, error) {
	return nil, i.executeBlock(b.Statements, environment.New(i.env))
}

func (i *Interpreter) VisitForVar(v StmtVar) (any, error) {
	var value any
	if v.Initializer != nil {
		v, err := i.evaluate(v.Initializer)
		if err != nil {
			return nil, err
		}
		value = v
	}
	i.env.Define(v.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitForVariable(v Variable) (any, error) {
	return i.env.Get(v.Name)
}

func (i *Interpreter) VisitForAssign(a Assign) (any, error) {
	value, err := i.evaluate(a.Value)
	if err != nil {
		return nil, err
	}
	if err := i.env.Assign(a.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitForGrouping(grouping Grouping) (any, error) {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) VisitForLiteral(literal Literal) (any, error) {
	return literal.Value, nil
}

func (i *Interpreter) VisitForBinary(binary Binary) (any, error) {
	left, err := i.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(binary.Right)
	if err != nil {
		return nil, err
	}

	switch binary.Operator.TokenType {
	case tokens.Plus:
		return sum(binary.Operator, left, right)
	case tokens.Minus:
		return numOperation(binary.Operator, left, right, func(l, r float64) any { return l - r })
	case tokens.Slash:
		return divide(binary.Operator, left, right)
	case tokens.Star:
		return numOperation(binary.Operator, left, right, func(l, r float64) any { return l * r })

	case tokens.Greater:
		return numOperation(binary.Operator, left, right, func(l, r float64) any { return l > r })
	case tokens.GreaterEqual:
		return numOperation(binary.Operator, left, right, func(l, r float64) any { return l >= r })
	case tokens.Less:
		return numOperation(binary.Operator, left, right, func(l, r float64) any { return l < r })
	case tokens.LessEqual:
		return numOperation(binary.Operator, left, right, func(l, r float64) any { return l <= r })

	case tokens.Equal: // Equality works the same as lox in go
		return left == right, nil
	case tokens.BangEqual:
		return left != right, nil

	}
	return nil, nil // unreachable
}

func (i *Interpreter) VisitForUnary(unary Unary) (any, error) {
	right, err := i.evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.TokenType {
	case tokens.Minus:
		r, err := asNumber(unary.Operator, right)
		if err != nil {
			return nil, err
		}
		return -r, nil
	case tokens.Bang:
		return !isTruthy(right), nil
	}
	return nil, nil // unreachable
}

func (i *Interpreter) evaluate(expression Expr) (any, error) {
	return expression.Accept(i)
}

// asNumber returns the number repesentation of the provided value or an error
// NOTE: it is Ok to try to parse as float because any number comes as a float
// due to the scanner implementation
func asNumber(op tokens.Token, v any) (float64, error) {
	f, ok := v.(float64)
	if !ok {
		return .0, errors.NewRuntimeError(op, "Operand must be a number.")
	}
	return f, nil
}

func asNumbers(op tokens.Token, a, b any) (float64, float64, error) {
	fa, ok := a.(float64)
	if !ok {
		return .0, .0, errors.NewRuntimeError(op, "Operand must be a number.")
	}
	fb, ok := b.(float64)
	if !ok {
		return .0, .0, errors.NewRuntimeError(op, "Operand must be a number.")
	}
	return fa, fb, nil
}

// numOperation executes the provided function with the arguments if they are numbers, returns an error otherwise
func numOperation(op tokens.Token, a, b any, f func(l, r float64) any) (any, error) {
	l, r, err := asNumbers(op, a, b)
	if err != nil {
		return nil, err
	}
	return f(l, r), nil
}

// isTruthy considers anything but nil or false value as true
func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	if value, isBool := v.(bool); isBool {
		return value
	}
	return true
}

// sum performs the '+' operation for either numbers or strings.
func sum(op tokens.Token, left any, right any) (any, error) {
	lNum, lIsNumber := left.(float64)
	rNum, rIsNumber := right.(float64)
	if lIsNumber && rIsNumber {
		return lNum + rNum, nil
	}
	lStr, lIsStr := left.(string)
	rStr, rIsStr := right.(string)
	if lIsStr && rIsStr {
		return lStr + rStr, nil
	}

	return nil, errors.NewRuntimeError(op, "Operands must be two numbers or two strings.")
}

func divide(op tokens.Token, left any, right any) (any, error) {
	l, r, err := asNumbers(op, left, right)
	if err != nil {
		return nil, err
	}
	if r == .0 {
		return nil, errors.NewRuntimeError(op, "Cannot divide by zero.")
	}
	return l / r, nil
}

// stringify returns the string representation of the provided value taking care of special cases for nil and numbers.
func stringify(v any) string {
	if v == nil {
		return "nil"
	}
	if f, ok := v.(float64); ok {
		strconv.FormatFloat(f, 'f', -1, 64)
	}
	return fmt.Sprintf("%v", v)
}
