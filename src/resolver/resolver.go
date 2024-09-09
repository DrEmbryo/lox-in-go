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

func (resolver *Resolver) Resolve(stmts []any) grammar.LoxError {
	for _, stmt := range stmts {
		switch entytyType := stmt.(type) {
		case grammar.Statement:
			resolver.resolveStmt(entytyType)
		case grammar.Expression:
			resolver.resolveExpr(entytyType)
		default:
			return nil
		}
	}
	return nil
}

func (resolver *Resolver) resolveStmt(stmt grammar.Statement) grammar.LoxError {
	switch stmtType := stmt.(type) {
	case grammar.VariableDeclarationStatement:
		return resolver.resolveVarStmt(stmtType)
	default:
		return nil
	}
}

func (resolver *Resolver) resolveVarStmt(stmt grammar.VariableDeclarationStatement) grammar.LoxError {
	err := resolver.declare(stmt.Name)
	if stmt.Initializer != nil {
		err := resolver.resolveExpr(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	resolver.define(stmt.Name)
	return err
}

func (resolver *Resolver) declare(name grammar.Token) grammar.LoxError {
	scope, err := resolver.scopes.Peek()
	if err != nil {
		return ResolverError{Token: name, Message: fmt.Sprint(err)}
	}
	scope[name.Lexeme.(string)] = false
	return nil
}

func (resolver *Resolver) define(name grammar.Token) grammar.LoxError {
	scope, err := resolver.scopes.Peek()
	if err != nil {
		return ResolverError{Token: name, Message: fmt.Sprint(err)}
	}
	scope[name.Lexeme.(string)] = true
	return nil
}

func (resolver *Resolver) resolveExpr(expr grammar.Expression) grammar.LoxError {
	switch exprType := expr.(type) {
	case grammar.VariableDeclaration:
		return resolver.resolveVarExpr(exprType)
	case grammar.AssignmentExpression:
		return resolver.resolveAssignmentExpr(exprType)
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

func (reolver *Resolver) resolveAssignmentExpr(expr grammar.AssignmentExpression) grammar.LoxError {
	err := reolver.resolveExpr(expr.Value)
	reolver.resolveLocal(expr, expr.Name)
	return err
}
