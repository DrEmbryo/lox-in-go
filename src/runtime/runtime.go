package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Runtime struct {
}

func (runtime *Runtime) literalExpr(expr grammar.LiteralExpression) any {
	return expr.Literal
}

func (runtime *Runtime) groupintExpr(expr  grammar.GroupingExpression) any {
	return runtime.evaluate(expr.Expression)
}

func (runtime *Runtime) unaryExpr(expr  grammar.UnaryExpression) any {
	right := runtime.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case  grammar.BANG:
		return !castToBool(right)
	case  grammar.MINUS:
		return right.(float64) * -1
	}

	return nil
}

func (runtime *Runtime) binaryExpr(expr  grammar.BinaryExpression) any {
	left := runtime.evaluate(expr.Left)
	right := runtime.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case  grammar.MINUS:
		return left.(float64) - right.(float64)
	case  grammar.PLUS:
		if checkTypeEquality(left, right) {
			switch left := left.(type) {
			case string:
				return fmt.Sprintf("%s%s", left, right)
			case float64:
				return left + right.(float64)
			}
		}
		return nil
	case  grammar.SLASH:
		return left.(float64) / right.(float64)
	case  grammar.STAR:
		return left.(float64) * right.(float64)
	case  grammar.GREATER:
		return left.(float64) > right.(float64)
	case  grammar.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case  grammar.LESS:
		return left.(float64) > right.(float64)
	case  grammar.LESS_EQUAL:
		return left.(float64) >= right.(float64)
	case  grammar.BANG_EQUAL:
		return !checkValueEquality(left, right)
	case  grammar.EQUAL_EQUAL:
		return checkValueEquality(left, right)
	}

	return nil
}

func (runtime *Runtime) evaluate(expr  grammar.Expression) any {

	switch exprType := expr.(type) {
	case  grammar.LiteralExpression:
		return runtime.literalExpr(exprType)
	case  grammar.GroupingExpression:
		return runtime.groupintExpr(exprType)
	case  grammar.UnaryExpression:
		return runtime.unaryExpr(exprType)
	case  grammar.BinaryExpression:
		return runtime.binaryExpr(exprType)
	}

	return nil
}

func checkTypeEquality(a, b any) bool {
	return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}

func checkValueEquality(a,b any) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func castToBool(val any) bool {
	switch v := val.(type) {
	case nil:
		return false
	case bool:
		return v
	}
	return true
}