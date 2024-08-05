package runtime

import (
	"fmt"
	"time"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Interpreter struct {
	Env       Environment
	globalEnv *Environment
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
	fmt.Printf("left: %v \n", left)
	if err != nil {
		return nil, err
	}
	right, err := interpreter.evaluate(expr.Right)
	fmt.Printf("right: %v \n", right)
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
			switch leftType := left.(type) {
			case string:
				return fmt.Sprintf("%s%s", left, right), nil
			case float64:
				return leftType + right.(float64), nil
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
		return left.(float64) < right.(float64), nil

	case grammar.LESS_EQUAL:
		err := checkNumericOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil

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

func (interpreter *Interpreter) callExpr(expr grammar.CallExpression) (any, grammar.LoxError) {
	var function LoxCallable
	callee, err := interpreter.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	switch calleeType := callee.(type) {
	case NativeCall:
		function = &calleeType
	default:
		return nil, RuntimeError{Token: expr.Paren, Message: "Calls available only for functions and classes"}
	}

	if len(expr.Arguments) != function.GetAirity() {
		return nil, RuntimeError{Token: expr.Paren, Message: fmt.Sprintf("Expect %v arguments but got %v.", function.GetAirity(), len(expr.Arguments))}
	}

	arguments := make([]any, 0)
	for argument := range expr.Arguments {
		arg, err := interpreter.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	return function.Call(*interpreter, arguments), nil
}

func (interpreter *Interpreter) evaluate(expr grammar.Expression) (any, grammar.LoxError) {
	switch exprType := expr.(type) {
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
	case grammar.CallExpression:
		return interpreter.callExpr(exprType)
	case grammar.LiteralExpression:
		return interpreter.literalExpr(exprType)
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

func (interpreter *Interpreter) expressionStmt(stmt grammar.ExpressionStatement) grammar.LoxError {
	_, err := interpreter.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	return nil
}

func (interpreter *Interpreter) whileStmt(stmt grammar.WhileLoopStatement) grammar.LoxError {
	expr, err := interpreter.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	for castToBool(expr) {
		expr, _ = interpreter.evaluate(stmt.Condition)
		_, err = interpreter.execute(stmt.Body)
		if err != nil {
			return err
		}
	}
	return err
}

func (interpreter *Interpreter) execute(stmt grammar.Statement) (any, grammar.LoxError) {
	switch stmtType := stmt.(type) {
	case grammar.PrintStatement:
		return nil, interpreter.printStmt(stmtType)
	case grammar.ExpressionStatement:
		return nil, interpreter.expressionStmt(stmtType)
	case grammar.VariableDeclarationStatement:
		return nil, interpreter.varStmt(stmtType)
	case grammar.BlockScopeStatement:
		return interpreter.blockStmt(stmtType)
	case grammar.ConditionalStatement:
		return nil, interpreter.conditionalStmt(stmtType)
	case grammar.WhileLoopStatement:
		return nil, interpreter.whileStmt(stmtType)
	case grammar.FunctionDeclarationStatement:
		return interpreter.functionDeclarationStmt(stmtType)
	case grammar.ReturnStatement:
		return interpreter.returnStmt(stmtType)
	}

	return nil, nil
}

func (interpreter *Interpreter) functionDeclarationStmt(stmt grammar.FunctionDeclarationStatement) (any, grammar.LoxError) {
	function := LoxFunction{Declaration: stmt}
	interpreter.Env.defineEnvValue(stmt.Name, function)
	return nil, nil
}

func (interpreter *Interpreter) returnStmt(stmt grammar.ReturnStatement) (any, grammar.LoxError) {
	return interpreter.evaluate(stmt.Expression)
}

func (interpreter *Interpreter) conditionalStmt(stmt grammar.ConditionalStatement) grammar.LoxError {
	condition, err := interpreter.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	if castToBool(condition) {
		_, err := interpreter.execute(stmt.ThenBranch)
		if err != nil {
			return err
		}
	} else if stmt.ElseBranch != nil {
		_, err := interpreter.execute(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return err
}

func (interpreter *Interpreter) blockStmt(stmt grammar.BlockScopeStatement) (any, grammar.LoxError) {
	parentEnv := interpreter.Env
	env := Environment{Values: make(map[string]any), Parent: &parentEnv}
	return interpreter.executeBlock(stmt.Statements, env)
}

func (interpreter *Interpreter) executeBlock(stmts []grammar.Statement, env Environment) (any, grammar.LoxError) {
	var err grammar.LoxError
	var value any
	parentEnv := interpreter.Env
	interpreter.Env = env
	for _, stmt := range stmts {
		value, err = interpreter.execute(stmt)
	}

	interpreter.Env = parentEnv
	return value, err
}

func (interpreter *Interpreter) Interpret(statements []grammar.Statement) []grammar.LoxError {
	interpreter.globalEnv = &interpreter.Env
	interpreter.globalEnv.defineEnvValue(grammar.Token{Lexeme: "clock"}, NativeCall{Airity: 0, NativeCallFunc: func(a ...any) any { return time.Now() }})

	errs := make([]grammar.LoxError, 0)
	for _, stmt := range statements {
		_, err := interpreter.execute(stmt)
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
	}
	interpreter.Env.defineEnvValue(stmt.Name, value)
	return err
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
