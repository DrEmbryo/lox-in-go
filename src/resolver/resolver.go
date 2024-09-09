// 183
package resolver

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/runtime"
	"github.com/DrEmbryo/lox/src/utils"
)

type Resolver struct {
	interpreter runtime.Interpreter
	scopes      utils.Stack[map[string]bool]
}

func (resolver *Resolver) beginScope() {
	resolver.scopes.Push(make(map[string]bool))
}

func (resolver *Resolver) endScope() {
	resolver.scopes.Pop()
}

func (resolver *Resolver) declare(name grammar.Token) {
	scope, err := resolver.scopes.Peek()
	if err != nil {
		return
	}
	scope[name.Lexeme.(string)] = false
}

func (resolver *Resolver) define(name grammar.Token) {
	scope, err := resolver.scopes.Peek()
	if err != nil {
		return
	}
	scope[name.Lexeme.(string)] = true
}

func (resolver *Resolver) Resolve(entity any) grammar.LoxError {
	fmt.Printf("%T", entity)
	switch entytyType := entity.(type) {
	case grammar.Statement:
		return resolver.resolveStmt(entytyType)
	case grammar.Expression:
		return resolver.resolveExpr(entytyType)
	default:
		return nil
	}
}

func (resolver *Resolver) resolveStmt(stmt grammar.Statement) grammar.LoxError {
	switch stmtType := stmt.(type) {
	case grammar.VariableDeclarationStatement:
		return resolver.resolveVarStmt(stmtType)
	case grammar.FunctionDeclarationStatement:
		return resolver.resolveFunctionStmt(stmtType)
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

func (resolver *Resolver) resolveVarStmt(stmt grammar.VariableDeclarationStatement) grammar.LoxError {
	resolver.declare(stmt.Name)
	if stmt.Initializer != nil {
		err := resolver.Resolve(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	resolver.define(stmt.Name)
	return nil
}

func (resolver *Resolver) resolveFunctionStmt(stmt grammar.FunctionDeclarationStatement) grammar.LoxError {
	resolver.declare(stmt.Name)
	resolver.define(stmt.Name)
	resolver.resolveFunction(stmt)
	return nil
}

func (resolver *Resolver) resolveFunction(function grammar.FunctionDeclarationStatement) grammar.LoxError {
	resolver.beginScope()
	for _, param := range function.Params {
		resolver.declare(param)
	}
	resolver.Resolve(function.Body)
	resolver.endScope()
	return nil
}

func (resolver *Resolver) resolveExpressionStmt(expr grammar.ExpressionStatement) grammar.LoxError {
	return resolver.Resolve(expr)
}

func (resolver *Resolver) resolveConditionalStmt(stmt grammar.ConditionalStatement) grammar.LoxError {
	err := resolver.Resolve(stmt.Condition)
	if err != nil {
		return err
	}
	err = resolver.Resolve(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		err = resolver.Resolve(stmt.ElseBranch)
	}
	return err
}

func (resolver *Resolver) resolvePrintStmt(stmt grammar.PrintStatement) grammar.LoxError {
	return resolver.Resolve(stmt.Value)
}

func (resolver *Resolver) resolveReturnStmt(stmt grammar.ReturnStatement) grammar.LoxError {
	if stmt.Expression != nil {
		return resolver.Resolve(stmt.Expression)
	}
	return nil
}

func (resolver *Resolver) resolveWhileStmt(stmt grammar.WhileLoopStatement) grammar.LoxError {
	err := resolver.Resolve(stmt.Condition)
	if err != nil {
		return err
	}
	return resolver.Resolve(stmt.Body)
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

func (resolver *Resolver) resolveVarExpr(expr grammar.VariableDeclaration) grammar.LoxError {
	scope, err := resolver.scopes.Peek()
	if err != nil {
		return ResolverError{Token: expr.Name, Message: fmt.Sprint(err)}
	}
	if val, ok := scope[expr.Name.Lexeme.(string)]; !resolver.scopes.IsEmpty() && ok && val {
		return ResolverError{Token: expr.Name, Message: "Can't read local variable in its own initializer."}
	}
	resolver.resolveLocal(expr, expr.Name)
	return nil
}

func (resolver *Resolver) resolveLocal(expr grammar.Expression, name grammar.Token) grammar.LoxError {
	for i := resolver.scopes.Len() - 1; i >= 0; i-- {
		scope, err := resolver.scopes.Get(i)
		if err != nil {
			return ResolverError{Token: name, Message: fmt.Sprint(err)}
		}
		if _, ok := scope[name.Lexeme.(string)]; ok {
			resolver.interpreter.Resolve(expr, resolver.scopes.Len()-1-i) //implement later
			return nil
		}
	}
	return nil
}

func (resolver *Resolver) resolveAssignmentExpr(expr grammar.AssignmentExpression) grammar.LoxError {
	err := resolver.Resolve(expr.Value)
	resolver.resolveLocal(expr, expr.Name)
	return err
}

func (resolver *Resolver) resolveBinaryExpr(expr grammar.BinaryExpression) grammar.LoxError {
	err := resolver.Resolve(expr.Left)
	if err != nil {
		return err
	}
	err = resolver.Resolve(expr.Right)
	return err
}

func (resolver *Resolver) resolveCallExpr(expr grammar.CallExpression) grammar.LoxError {
	err := resolver.Resolve(expr.Callee)
	for _, argument := range expr.Arguments {
		err := resolver.Resolve(argument)
		if err != nil {
			return err
		}
	}
	return err
}

func (resolver *Resolver) resolveGroupExpr(expr grammar.GroupingExpression) grammar.LoxError {
	return resolver.Resolve(expr.Expression)
}

func (resolver *Resolver) resolveLiteralExpr() grammar.LoxError {
	return nil
}

func (resolver *Resolver) resolveUnaryExpr(expr grammar.UnaryExpression) grammar.LoxError {
	return resolver.Resolve(expr.Right)
}
