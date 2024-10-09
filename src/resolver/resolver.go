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
	METHOD
	CLASS
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
	key := name.Lexeme.(string)
	_, ok := scope[key]
	if ok {
		return ResolverError{Token: name, Message: "Already variable with this name in this scope."}
	} else {
		scope[key] = false
	}
	return nil
}

func (resolver *Resolver) define(name grammar.Token) {
	scope, err := resolver.Scopes.Peek()
	if err != nil {
		return
	}
	key := name.Lexeme.(string)
	scope[key] = true
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
		return resolver.resolveBlockStmt(stmtType)
	case grammar.VariableDeclarationStatement:
		return resolver.resolveVarStmt(stmtType)
	case grammar.FunctionDeclarationStatement:
		return resolver.resolveFunctionStmt(stmtType)
	case grammar.ClassDeclarationStatement:
		return resolver.resolveClassStmt(stmtType)
	case grammar.ExpressionStatement:
		return resolver.resolveExpressionStmt(stmtType)
	case grammar.ConditionalStatement:
		return resolver.resolveConditionalStmt(stmtType)
	case grammar.PrintStatement:
		return resolver.resolvePrintStmt(stmtType)
	case grammar.ReturnStatement:
		return resolver.resolveReturnStmt(stmtType)
	case grammar.WhileLoopStatement:
		return resolver.resolveWhileStmt(stmtType)
	default:
		return nil
	}
}

func (resolver *Resolver) resolveBlockStmt(stmt grammar.BlockScopeStatement) grammar.LoxError {
	resolver.beginScope()
	resolver.resolveStmts(stmt.Statements)
	resolver.endScope()
	return nil
}

func (resolver *Resolver) resolveVarStmt(stmt grammar.VariableDeclarationStatement) grammar.LoxError {
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

func (resolver *Resolver) resolveFunctionStmt(stmt grammar.FunctionDeclarationStatement) grammar.LoxError {
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

func (resolver *Resolver) resolveClassStmt(class grammar.ClassDeclarationStatement) grammar.LoxError {
	enclosingClass := resolver.CurrentClass
	resolver.CurrentClass = CLASS
	err := resolver.declare(class.Name)
	if err != nil {
		return err
	}
	resolver.define(class.Name)
	resolver.beginScope()
	scope, stackErr := resolver.Scopes.Peek()
	if stackErr != nil {
		return ResolverError{Token: class.Name, Message: fmt.Sprint(stackErr)}
	}
	scope["this"] = true
	for _, method := range class.Methods {
		declaration := METHOD
		resolver.resolveFunction(method, declaration)
	}
	resolver.endScope()
	resolver.CurrentClass = enclosingClass
	return err
}

func (resolver *Resolver) resolveExpressionStmt(expr grammar.ExpressionStatement) grammar.LoxError {
	return resolver.resolveExpr(expr)
}

func (resolver *Resolver) resolveConditionalStmt(stmt grammar.ConditionalStatement) grammar.LoxError {
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

func (resolver *Resolver) resolvePrintStmt(stmt grammar.PrintStatement) grammar.LoxError {
	return resolver.resolveExpr(stmt.Value)
}

func (resolver *Resolver) resolveReturnStmt(stmt grammar.ReturnStatement) grammar.LoxError {
	if resolver.CurrentFunction == NONE {
		return ResolverError{Token: stmt.Keyword, Message: "Can't return from top-level code."}
	}

	if stmt.Expression != nil {
		return resolver.resolveExpr(stmt.Expression)
	}
	return nil
}

func (resolver *Resolver) resolveWhileStmt(stmt grammar.WhileLoopStatement) grammar.LoxError {
	err := resolver.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	return resolver.resolveStmt(stmt.Body)
}

func (resolver *Resolver) resolveExpr(expr grammar.Expression) grammar.LoxError {
	switch exprType := expr.(type) {
	case grammar.VariableDeclaration:
		return resolver.resolveVarExpr(exprType)
	case grammar.AssignmentExpression:
		return resolver.resolveAssignmentExpr(exprType)
	case grammar.BinaryExpression:
		return resolver.resolveBinaryExpr(exprType)
	case grammar.CallExpression:
		return resolver.resolveCallExpr(exprType)
	case grammar.PropertyAccessExpression:
		return resolver.propAccessExpr(exprType)
	case grammar.PropertyAssignmentExpression:
		return resolver.propAssignmentExpr(exprType)
	case grammar.SelfReferenceExpression:
		return resolver.selfReferenceExpr(exprType)
	case grammar.GroupingExpression:
		return resolver.resolveGroupExpr(exprType)
	case grammar.LiteralExpression:
		return resolver.resolveLiteralExpr()
	case grammar.UnaryExpression:
		return resolver.resolveUnaryExpr(exprType)
	default:
		return nil
	}
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

func (resolver *Resolver) resolveVarExpr(expr grammar.VariableDeclaration) grammar.LoxError {
	scope, err := resolver.Scopes.Peek()
	if err != nil {
		return ResolverError{Token: expr.Name, Message: fmt.Sprint(err)}
	}
	key := expr.Name.Lexeme.(string)
	if val, ok := scope[key]; !resolver.Scopes.IsEmpty() && ok && !val {
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
		key := name.Lexeme.(string)
		if _, ok := scope[key]; ok {
			resolver.Interpreter.Resolve(expr, resolver.Scopes.Len()-1-i)
			return nil
		}
	}
	return nil
}

func (resolver *Resolver) resolveAssignmentExpr(expr grammar.AssignmentExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Value)
	resolver.resolveLocal(expr, expr.Name)
	return err
}

func (resolver *Resolver) resolveBinaryExpr(expr grammar.BinaryExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Left)
	if err != nil {
		return err
	}
	err = resolver.resolveExpr(expr.Right)
	return err
}

func (resolver *Resolver) resolveCallExpr(expr grammar.CallExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		err := resolver.resolveExpr(argument)
		if err != nil {
			return err
		}
	}
	return err
}

func (resolver *Resolver) resolveGroupExpr(expr grammar.GroupingExpression) grammar.LoxError {
	return resolver.resolveExpr(expr.Expression)
}

func (resolver *Resolver) resolveLiteralExpr() grammar.LoxError {
	return nil
}

func (resolver *Resolver) resolveUnaryExpr(expr grammar.UnaryExpression) grammar.LoxError {
	return resolver.resolveExpr(expr.Right)
}
