package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Interpriter struct {
	Env Environment
}

func (interpriter *Interpriter) literalExpr(expr grammar.LiteralExpression) (any, grammar.LoxError) {
	return expr.Literal, nil
}

func (interpriter *Interpriter) groupingExpr(expr  grammar.GroupingExpression) (any, grammar.LoxError) {
	return interpriter.evaluate(expr.Expression)
}

func (interpriter *Interpriter) unaryExpr(expr  grammar.UnaryExpression) (any, grammar.LoxError) {
	right, err := interpriter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case  grammar.BANG:
		return !castToBool(right), nil
	case  grammar.MINUS:
		err := checkNuericOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return right.(float64) * -1, nil
	}
	return nil, nil
}

func (interpriter *Interpriter) binaryExpr(expr  grammar.BinaryExpression) (any, grammar.LoxError) {
	left, err := interpriter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := interpriter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case  grammar.MINUS:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil

	case  grammar.PLUS:
		if checkTypeEquality(left, right) {
			switch left := left.(type) {
			case string:
				return fmt.Sprintf("%s%s", left, right), nil
			case float64:
				return left + right.(float64), nil
			}
		}
		return nil, RuntimeError{Token: expr.Operator, Message: "Operands must be two numbers or two strings."}

	case  grammar.SLASH:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil

	case  grammar.STAR:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	case  grammar.GREATER:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case  grammar.GREATER_EQUAL:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil

	case  grammar.LESS:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case  grammar.LESS_EQUAL:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil

	case  grammar.BANG_EQUAL:
		return !checkValueEquality(left, right), nil
	case  grammar.EQUAL_EQUAL:
		return checkValueEquality(left, right), nil
	}

	return nil, nil
}

func (interpriter *Interpriter) evaluate(expr  grammar.Expression) (any, grammar.LoxError) {
	switch exprType := expr.(type) {
	case  grammar.LiteralExpression:
		return interpriter.literalExpr(exprType)
	case  grammar.GroupingExpression:
		return interpriter.groupingExpr(exprType)
	case  grammar.UnaryExpression:
		return interpriter.unaryExpr(exprType)
	case  grammar.BinaryExpression:
		return interpriter.binaryExpr(exprType)
	case grammar.VariableDeclaration:
		return interpriter.varExpr(exprType)
	case grammar.AssignmentExpression:
		return interpriter.assignmentExpr(exprType)
	default:
		fmt.Printf("%T", exprType)
	}

	return nil, nil
}

func (interpriter *Interpriter) printStmt(stmt grammar.PrintStatment) grammar.LoxError {
	value, err := interpriter.evaluate(stmt.Value)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (interpriter *Interpriter) expressionStmt(stmt grammar.Statement) grammar.LoxError {
	_, err := interpriter.evaluate(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (interpriter *Interpriter) execute(stmt grammar.Statement) grammar.LoxError {
	switch stmtType := stmt.(type) {
	case  grammar.PrintStatment:
		return interpriter.printStmt(stmtType)
	case  grammar.ExpressionStatement:
		return interpriter.expressionStmt(stmtType)
	case grammar.VariableDeclarationStatment:
		return interpriter.varStmt(stmtType)
	}

	return nil
}

func (interpriter Interpriter) Interpret(statements []grammar.Statement) []grammar.LoxError {
	errs := make([]grammar.LoxError, 0)
	for _, stmt := range statements {
		err := interpriter.execute(stmt)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return  errs
}

func (interpriter *Interpriter) varStmt(stmt grammar.VariableDeclarationStatment) grammar.LoxError {
	var value any
	var err grammar.LoxError
	if stmt.Initializer != nil {
		value, err = interpriter.evaluate(stmt.Initializer)
		interpriter.Env.defineEnvValue(fmt.Sprintf("%s", stmt.Name.Lexeme), value)
		return err
	} 
	return RuntimeError{Token: stmt.Name, Message: "Expect initialization of variable"}
}

func (interpriter *Interpriter) varExpr(expr grammar.VariableDeclaration) (any, grammar.LoxError) {
	return interpriter.Env.getEnvValue(expr.Name) 
}

func (interpriter *Interpriter) assignmentExpr(expr grammar.AssignmentExpression) (any, grammar.LoxError) {
	value, err := interpriter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	interpriter.Env.assignEnvValue(expr.Name, value)
	return value, nil
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

func checkNuericOperand(operator grammar.Token, operand any) grammar.LoxError {
	switch operand.(type) {
	case float64:
		return nil
	}
	return RuntimeError{Token: operator, Message: "Operand must be a number."}
}

func checkNumericOperands(operator grammar.Token, left any, right any) grammar.LoxError {
	if checkTypeEquality(left, right) {
		switch left.(type) {
		case float64:
			return nil
		}
	}
	return RuntimeError{Token: operator, Message: "Operands must be numbers."}
}