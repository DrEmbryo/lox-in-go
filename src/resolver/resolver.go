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
	case grammar.FunctionDeclarationStatement:
		return resolver.resolveFunctionStmt(stmtType)
	default:
		return nil
	}
}

func (resolver *Resolver) resolveVarStmt(stmt grammar.VariableDeclarationStatement) grammar.LoxError {
	resolver.declare(stmt.Name)
	if stmt.Initializer != nil {
		err := resolver.resolveExpr(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	resolver.define(stmt.Name)
	return nil
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

func (resolver *Resolver) resolveAssignmentExpr(expr grammar.AssignmentExpression) grammar.LoxError {
	err := resolver.resolveExpr(expr.Value)
	resolver.resolveLocal(expr, expr.Name)
	return err
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
	resolver.resolveStmt(function.Body)
	resolver.endScope()
	return nil
}
