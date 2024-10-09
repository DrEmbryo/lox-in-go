package runtime

import (
	"fmt"
	"time"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Interpreter struct {
	Env       Environment
	globalEnv *Environment
	LocalEnv  map[any]int
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
	return nil, err
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
	return nil, err
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
	case LoxFunction:
		function = &calleeType
	case NativeCall:
		function = &calleeType
	case LoxClass:
		function = &calleeType
	default:
		return nil, RuntimeError{Token: expr.Paren, Message: "Calls available only for functions and classes"}
	}

	if len(expr.Arguments) != function.GetAirity() {
		return nil, RuntimeError{Token: expr.Paren, Message: fmt.Sprintf("Expect %v arguments but got %v.", function.GetAirity(), len(expr.Arguments))}
	}

	arguments := make([]any, 0)
	for _, argument := range expr.Arguments {
		arg, err := interpreter.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	return function.Call(*interpreter, arguments)
}

func (interpreter *Interpreter) propAccessExpr(expr grammar.PropertyAccessExpression) (any, grammar.LoxError) {
	object, err := interpreter.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	classInstance, ok := object.(LoxClassInstance)
	if ok {
		return classInstance.GetProperty(expr.Name)
	}

	return nil, RuntimeError{Token: expr.Name, Message: "Only instances have prooperties."}
}

func (interpreter *Interpreter) propAssignmentExpr(expr grammar.PropertyAssignmentExpression) (any, grammar.LoxError) {
	object, err := interpreter.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if _, ok := object.(LoxClassInstance); !ok {
		return nil, RuntimeError{Token: expr.Name, Message: "Only instances have fiellds."}
	}

	value, err := interpreter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if classInstance, ok := object.(LoxClassInstance); ok {
		return classInstance.SetProperty(expr.Name, value), nil
	}
	return value, err
}

func (interpreter *Interpreter) selfReferenceExpr(expr grammar.SelfReferenceExpression) (any, grammar.LoxError) {
	return interpreter.lookUpVariable(expr.Keyword, expr)
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
	case grammar.PropertyAccessExpression:
		return interpreter.propAccessExpr(exprType)
	case grammar.PropertyAssignmentExpression:
		return interpreter.propAssignmentExpr(exprType)
	case grammar.SelfReferenceExpression:
		return interpreter.selfReferenceExpr(exprType)
	case grammar.LiteralExpression:
		return interpreter.literalExpr(exprType)
	default:
		fmt.Printf("%T", exprType)
		return nil, nil
	}
}

func (interpreter *Interpreter) printStmt(stmt grammar.PrintStatement) grammar.LoxError {
	value, err := interpreter.evaluate(stmt.Value)
	if err != nil {
		return err
	}

	switch valueType := value.(type) {
	case LoxClassInstance:
		fmt.Println(valueType.ToString())
	case LoxClass:
		fmt.Println(valueType.ToString())
	case LoxFunction:
		fmt.Println(valueType.ToString())
	default:
		fmt.Println(valueType)
	}

	return err
}

func (interpreter *Interpreter) expressionStmt(stmt grammar.ExpressionStatement) grammar.LoxError {
	_, err := interpreter.evaluate(stmt.Expression)
	return err
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
	case grammar.ClassDeclarationStatement:
		return interpreter.classDeclarationStmt(stmtType)
	case grammar.ReturnStatement:
		return interpreter.returnStmt(stmtType)
	default:
		return nil, nil
	}
}

func (interpreter *Interpreter) functionDeclarationStmt(stmt grammar.FunctionDeclarationStatement) (any, grammar.LoxError) {
	function := LoxFunction{Declaration: stmt, Closure: interpreter.Env}
	interpreter.Env.defineEnvValue(stmt.Name, function)
	return nil, nil
}

func (interpreter *Interpreter) classDeclarationStmt(stmt grammar.ClassDeclarationStatement) (any, grammar.LoxError) {
	methods := make(map[string]LoxFunction)

	for _, method := range stmt.Methods {
		lookup := fmt.Sprintf("%s", method.Name.Lexeme)
		methods[lookup] = LoxFunction{Closure: interpreter.Env, Declaration: method}
	}

	class := LoxClass{Name: stmt.Name, Methods: methods, Fields: make(map[any]any)}
	interpreter.Env.defineEnvValue(stmt.Name, class)
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
	return interpreter.lookUpVariable(expr.Name, expr)
}

func (interpreter *Interpreter) lookUpVariable(name grammar.Token, expr grammar.Expression) (any, grammar.LoxError) {
	distance, ok := interpreter.LocalEnv[expr]
	if !ok {
		return interpreter.Env.getEnvValueAt(distance, name)
	}
	return interpreter.globalEnv.getEnvValue(name)
}

func (interpreter *Interpreter) assignmentExpr(expr grammar.AssignmentExpression) (any, grammar.LoxError) {
	value, err := interpreter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	distance, ok := interpreter.LocalEnv[expr]
	if !ok {
		interpreter.Env.assignEnvValueAt(distance, expr.Name, value)
	} else {
		interpreter.globalEnv.assignEnvValue(expr.Name, value)
	}
	return value, err
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

func (interpreter *Interpreter) Resolve(expr grammar.Expression, depth int) {
	interpreter.LocalEnv[expr] = depth
}
