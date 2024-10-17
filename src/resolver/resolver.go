package resolver

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/runtime"
	"github.com/DrEmbryo/lox/src/utils"
)

const (
	NONE = iota
	FUNCTION
	CLASS
	SUBCLASS
	METHOD
	INITIALIZER
)

type Resolver struct {
	Interpreter     runtime.Interpreter
	Scopes          utils.Stack[map[string]bool]
	Error           []grammar.LoxError
	CurrentFunction int
	CurrentClass    int
}

func (resolver *Resolver) beginScope() {
	resolver.Scopes.Push(make(map[string]bool, 0))
}

func (resolver *Resolver) endScope() {
	resolver.Scopes.Pop()
}

func (resolver *Resolver) declare(name grammar.Token) grammar.LoxError {
	scope, err := resolver.Scopes.Peek()
	if err != nil {
		return nil
	}
	lookup := fmt.Sprintf("%s", name.Lexeme)
	_, ok := scope[lookup]
	if ok {
		return ResolverError{Token: name, Message: "Already variable with this name in this scope."}
	} else {
		scope[lookup] = false
	}
	return nil
}

func (resolver *Resolver) define(name grammar.Token) {
	scope, err := resolver.Scopes.Peek()
	if err != nil {
		return
	}
	lookup := fmt.Sprintf("%s", name.Lexeme)
	scope[lookup] = true
}

func (resolver *Resolver) Resolve(statements []grammar.Statement) []grammar.LoxError {
	resolver.beginScope()
	resolver.resolveStmts(statements)
	resolver.endScope()
	return resolver.Error
}

func (resolver *Resolver) resolveStmts(statements []grammar.Statement) {
	for _, stmt := range statements {
		err := resolver.resolveStmt(stmt)
		if err != nil {
			resolver.Error = append(resolver.Error, err)
		}
	}
}

func (resolver *Resolver) resolveStmt(stmt grammar.Statement) grammar.LoxError {
	switch stmtType := stmt.(type) {
	case grammar.BlockScopeStatement:
		return resolver.blockStmt(stmtType)
	case grammar.VariableDeclarationStatement:
		return resolver.vararStmt(stmtType)
	case grammar.FunctionDeclarationStatement:
		return resolver.functionStmt(stmtType)
	case grammar.ClassDeclarationStatement:
		return resolver.classStmt(stmtType)
	case grammar.ExpressionStatement:
		return resolver.expressionStmt(stmtType)
	case grammar.ConditionalStatement:
		return resolver.conditionalStmt(stmtType)
	case grammar.PrintStatement:
		return resolver.printStmt(stmtType)
	case grammar.ReturnStatement:
		return resolver.returnStmt(stmtType)
	case grammar.WhileLoopStatement:
		return resolver.whileStmt(stmtType)
	default:
		return nil
	}
}

func (resolver *Resolver) blockStmt(stmt grammar.BlockScopeStatement) grammar.LoxError {
	resolver.beginScope()
	resolver.resolveStmts(stmt.Statements)
	resolver.endScope()
	return nil
}

func (resolver *Resolver) vararStmt(stmt grammar.VariableDeclarationStatement) grammar.LoxError {
	err := resolver.declare(stmt.Name)
	if err != nil {
		return err
	}
	if stmt.Initializer != nil {
		err := resolver.resolveStmt(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	resolver.define(stmt.Name)
	return err
}

func (resolver *Resolver) functionStmt(stmt grammar.FunctionDeclarationStatement) grammar.LoxError {
	err := resolver.declare(stmt.Name)
	resolver.define(stmt.Name)
	resolver.resolveFunction(stmt, FUNCTION)
	return err
}

func (resolver *Resolver) resolveFunction(function grammar.FunctionDeclarationStatement, functionType int) grammar.LoxError {
	enclosingFunction := resolver.CurrentFunction
	resolver.CurrentFunction = functionType
	resolver.beginScope()
	for _, param := range function.Params {
		err := resolver.declare(param)
		if err != nil {
			return err
		}
		resolver.define(param)
	}
	resolver.resolveStmt(function.Body)
	resolver.endScope()
	resolver.CurrentFunction = enclosingFunction
	return nil
}

func (resolver *Resolver) classStmt(class grammar.ClassDeclarationStatement) grammar.LoxError {
	enclosingClass := resolver.CurrentClass
	resolver.CurrentClass = CLASS
	resolver.declare(class.Name)
	resolver.define(class.Name)
	super, ok := class.Super.(grammar.VariableDeclaration)
	if ok {
		if class.Name.Lexeme == super.Name.Lexeme {
			return ResolverError{Token: super.Name, Message: "A class can't inherit from itself."}
		}
		resolver.CurrentClass = SUBCLASS
		resolver.resolveExpr(super)
	}
	if class.Super != nil {
		resolver.beginScope()
		superScope, stackErr := resolver.Scopes.Peek()
		if stackErr != nil {
			return ResolverError{Token: class.Name, Message: fmt.Sprint(stackErr)}
		}
		superScope["super"] = true
	}
	resolver.beginScope()
	scope, stackErr := resolver.Scopes.Peek()
	if stackErr != nil {
		return ResolverError{Token: class.Name, Message: fmt.Sprint(stackErr)}
	}
	scope["this"] = true
	for _, method := range class.Methods {
		declaration := METHOD
		if method.Name.Lexeme == runtime.CONSTRUCTOR {
			declaration = INITIALIZER
		}
		resolver.resolveFunction(method, declaration)
	}
	if class.Super != nil {
		resolver.endScope()
	}
	resolver.endScope()
	resolver.CurrentClass = enclosingClass
	return nil
}

func (resolver *Resolver) expressionStmt(expr grammar.ExpressionStatement) grammar.LoxError {
	return resolver.resolveExpr(expr)
}

func (resolver *Resolver) conditionalStmt(stmt grammar.ConditionalStatement) grammar.LoxError {
	err := resolver.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = resolver.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		err = resolver.resolveStmt(stmt.ElseBranch)
	}
	return err
}

func (resolver *Resolver) printStmt(stmt grammar.PrintStatement) grammar.LoxError {
	return resolver.resolveExpr(stmt.Value)
}

func (resolver *Resolver) returnStmt(stmt grammar.ReturnStatement) grammar.LoxError {
	if resolver.CurrentFunction == NONE {
		return ResolverError{Token: stmt.Keyword, Message: "Can't return from top-level code."}
	}

	if stmt.Expression != nil {
		if resolver.CurrentFunction == INITIALIZER {
			return ResolverError{Token: stmt.Keyword, Message: "Can't return a value from constructor"}
		}
		return resolver.resolveExpr(stmt.Expression)
	}
	return nil
}

func (resolver *Resolver) whileStmt(stmt grammar.WhileLoopStatement) grammar.LoxError {
	err := resolver.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	return resolver.resolveStmt(stmt.Body)
}

func (resolver *Resolver) resolveExpr(expr grammar.Expression) grammar.LoxError {
	switch exprType := expr.(type) {
	case grammar.VariableDeclaration:
		return resolver.varExpr(exprType)
	case grammar.AssignmentExpression:
		return resolver.assignmentExpr(exprType)
	case grammar.BinaryExpression:
		return resolver.binaryExpr(exprType)
	case grammar.CallExpression:
		return resolver.callExpr(exprType)
	case grammar.PropertyAccessExpression:
		return resolver.propAccessExpr(exprType)
	case grammar.PropertyAssignmentExpression:
		return resolver.propAssignmentExpr(exprType)
	case grammar.SelfReferenceExpression:
		return resolver.selfReferenceExpr(exprType)
	case grammar.BaseClassCallExpression:
		return resolver.baseClassCallExpr(exprType)
	case grammar.GroupingExpression:
		return resolver.groupExpr(exprType)
	case grammar.LiteralExpression:
		return resolver.literalExpr()
	case grammar.UnaryExpression:
		return resolver.unaryExpr(exprType)
	default:
		return nil
	}
}

func (resolver *Resolver) baseClassCallExpr(expr grammar.BaseClassCallExpression) grammar.LoxError {
	if resolver.CurrentClass == NONE {
		return ResolverError{Token: expr.Keyword, Message: "Can't use 'super' outside of a class."}
	} else if resolver.CurrentClass != SUBCLASS {
		return ResolverError{Token: expr.Keyword, Message: "Can't use 'super' in a class with no superclass."}
	}
	return resolver.resolveLocal(expr, expr.Keyword)
}

func (resolver *Resolver) selfReferenceExpr(expr grammar.SelfReferenceExpression) grammar.LoxError {
	if resolver.CurrentClass == NONE {
		return runtime.RuntimeError{Token: expr.Keyword, Message: "Can't use 'this' outside of a class"}
	}
	return resolver.resolveLocal(expr, expr.Keyword)
}

func (resolver *Resolver) propAssignmentExpr(expr grammar.PropertyAssignmentExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Value)
	if err != nil {
		return err
	}
	err = resolver.resolveExpr(expr.Object)
	return err
}

func (resolver *Resolver) propAccessExpr(expr grammar.PropertyAccessExpression) grammar.LoxError {
	return resolver.resolveExpr(expr.Object)
}

func (resolver *Resolver) varExpr(expr grammar.VariableDeclaration) grammar.LoxError {
	scope, err := resolver.Scopes.Peek()
	if err != nil {
		return ResolverError{Token: expr.Name, Message: fmt.Sprint(err)}
	}
	lookup := fmt.Sprintf("%s", expr.Name.Lexeme)
	if val, ok := scope[lookup]; !resolver.Scopes.IsEmpty() && ok && !val {
		return ResolverError{Token: expr.Name, Message: "Can't read local variable in its own initializer."}
	}
	resolver.resolveLocal(expr, expr.Name)
	return nil
}

func (resolver *Resolver) resolveLocal(expr grammar.Expression, name grammar.Token) grammar.LoxError {
	for i := resolver.Scopes.Len() - 1; i >= 0; i-- {
		scope, err := resolver.Scopes.Get(i)
		if err != nil {
			return ResolverError{Token: name, Message: fmt.Sprint(err)}
		}
		lookup := fmt.Sprintf("%s", name.Lexeme)
		if _, ok := scope[lookup]; ok {
			resolver.Interpreter.Resolve(expr, resolver.Scopes.Len()-1-i)
			return nil
		}
	}
	return nil
}

func (resolver *Resolver) assignmentExpr(expr grammar.AssignmentExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Value)
	resolver.resolveLocal(expr, expr.Name)
	return err
}

func (resolver *Resolver) binaryExpr(expr grammar.BinaryExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Left)
	if err != nil {
		return err
	}
	err = resolver.resolveExpr(expr.Right)
	return err
}

func (resolver *Resolver) callExpr(expr grammar.CallExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		err := resolver.resolveExpr(argument)
		if err != nil {
			return err
		}
	}
	return err
}

func (resolver *Resolver) groupExpr(expr grammar.GroupingExpression) grammar.LoxError {
	return resolver.resolveExpr(expr.Expression)
}

func (resolver *Resolver) literalExpr() grammar.LoxError {
	return nil
}

func (resolver *Resolver) unaryExpr(expr grammar.UnaryExpression) grammar.LoxError {
	return resolver.resolveExpr(expr.Right)
}
