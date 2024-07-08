package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Interpreter struct {
	Env Environment
}

func (interpreter *Interpreter) literalExpr(expr grammar.LiteralExpression) (any, grammar.LoxError) {
	return expr.Literal, nil
}

func (interpreter *Interpreter) groupingExpr(expr grammar.GroupingExpression) (any, grammar.LoxError) {
	return interpreter.evaluate(expr.Expression)
}

func (interpreter *Interpreter) unaryExpr(expr grammar.UnaryExpression) (any, grammar.LoxError) {
	right, err := interpreter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case grammar.BANG:
		return !castToBool(right), nil
	case grammar.MINUS:
		err := checkNumericOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return right.(float64) * -1, nil
	}
	return nil, nil
}

func (interpreter *Interpreter) binaryExpr(expr grammar.BinaryExpression) (any, grammar.LoxError) {
	left, err := interpreter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := interpreter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case grammar.MINUS:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil

	case grammar.PLUS:
		if checkTypeEquality(left, right) {
			switch left := left.(type) {
			case string:
				return fmt.Sprintf("%s%s", left, right), nil
			case float64:
				return left + right.(float64), nil
			}
		}
		return nil, RuntimeError{Token: expr.Operator, Message: "Operands must be two numbers or two strings."}

	case grammar.SLASH:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil

	case grammar.STAR:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	case grammar.GREATER:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case grammar.GREATER_EQUAL:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil

	case grammar.LESS:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case grammar.LESS_EQUAL:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil

	case grammar.BANG_EQUAL:
		return !checkValueEquality(left, right), nil
	case grammar.EQUAL_EQUAL:
		return checkValueEquality(left, right), nil
	}

	return nil, nil
}

func (interpreter *Interpreter) logicalExpr(expr grammar.LogicExpression) (any, grammar.LoxError) {
	left, err := interpreter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == grammar.OR {
		if castToBool(left) {
			return left, nil
		} else if !castToBool(left) {
			return left, nil
		}
	}

	return interpreter.evaluate(expr.Right)
}

func (interpreter *Interpreter) evaluate(expr grammar.Expression) (any, grammar.LoxError) {
	switch exprType := expr.(type) {
	case grammar.LiteralExpression:
		return interpreter.literalExpr(exprType)
	case grammar.GroupingExpression:
		return interpreter.groupingExpr(exprType)
	case grammar.UnaryExpression:
		return interpreter.unaryExpr(exprType)
	case grammar.BinaryExpression:
		return interpreter.binaryExpr(exprType)
	case grammar.VariableDeclaration:
		return interpreter.varExpr(exprType)
	case grammar.AssignmentExpression:
		return interpreter.assignmentExpr(exprType)
	case grammar.LogicExpression:
		return interpreter.logicalExpr(exprType)
	default:
		fmt.Printf("%T", exprType)
	}

	return nil, nil
}

func (interpreter *Interpreter) printStmt(stmt grammar.PrintStatement) grammar.LoxError {
	value, err := interpreter.evaluate(stmt.Value)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (interpreter *Interpreter) expressionStmt(stmt grammar.Statement) grammar.LoxError {
	_, err := interpreter.evaluate(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (interpreter *Interpreter) execute(stmt grammar.Statement) grammar.LoxError {
	switch stmtType := stmt.(type) {
	case grammar.PrintStatement:
		return interpreter.printStmt(stmtType)
	case grammar.ExpressionStatement:
		return interpreter.expressionStmt(stmtType)
	case grammar.VariableDeclarationStatement:
		return interpreter.varStmt(stmtType)
	case grammar.BlockScopeStatement:
		return interpreter.blockStmt(stmtType)
	case grammar.ConditionalStatement:
		return interpreter.conditionalStmt(stmtType)
	}

	return nil
}

func (interpreter *Interpreter) conditionalStmt(stmt grammar.ConditionalStatement) grammar.LoxError {
	condition, err := interpreter.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	if castToBool(condition) {
		err := interpreter.execute(stmt.ThenBranch)
		if err != nil {
			return err
		}
	} else if stmt.ElseBranch != nil {
		err := interpreter.execute(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return err
}

func (interpreter *Interpreter) blockStmt(stmt grammar.BlockScopeStatement) grammar.LoxError {
	parentEnv := interpreter.Env
	env := Environment{Values: make(map[string]any), Parent: &parentEnv}
	return interpreter.executeBlock(stmt.Statements, env)
}

func (interpreter *Interpreter) executeBlock(stmts []grammar.Statement, env Environment) grammar.LoxError {
	var err grammar.LoxError
	parentEnv := interpreter.Env
	interpreter.Env = env
	for _, stmt := range stmts {
		err = interpreter.execute(stmt)
	}

	interpreter.Env = parentEnv
	return err
}

func (interpreter *Interpreter) Interpret(statements []grammar.Statement) []grammar.LoxError {
	errs := make([]grammar.LoxError, 0)
	for _, stmt := range statements {
		err := interpreter.execute(stmt)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (interpreter *Interpreter) varStmt(stmt grammar.VariableDeclarationStatement) grammar.LoxError {
	var value any
	var err grammar.LoxError
	if stmt.Initializer != nil {
		value, err = interpreter.evaluate(stmt.Initializer)
		interpreter.Env.defineEnvValue(fmt.Sprintf("%s", stmt.Name.Lexeme), value)
		return err
	}
	return RuntimeError{Token: stmt.Name, Message: "Expect initialization of variable"}
}

func (interpreter *Interpreter) varExpr(expr grammar.VariableDeclaration) (any, grammar.LoxError) {
	return interpreter.Env.getEnvValue(expr.Name)
}

func (interpreter *Interpreter) assignmentExpr(expr grammar.AssignmentExpression) (any, grammar.LoxError) {
	value, err := interpreter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	interpreter.Env.assignEnvValue(expr.Name, value)
	return value, nil
}

func checkTypeEquality(a, b any) bool {
	return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}

func checkValueEquality(a, b any) bool {
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

func checkNumericOperand(operator grammar.Token, operand any) grammar.LoxError {
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
