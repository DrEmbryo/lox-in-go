package utils

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/DrEmbryo/jlox/src/grammar"
)

type AstPrinter struct {
}

type TokenPrinter struct {
}

func (printer *AstPrinter) Print(stmts []grammar.Statement) {
	fmt.Println("Ast generated from tokens:")
	for _, stmt := range stmts {
		offset := 0
		fmt.Println(printer.printNode(offset, stmt))
	}
	fmt.Println("")
}

func (printer *AstPrinter) printNode(offset int, stmt grammar.Statement) string {
	nodeType := fmt.Sprintf("%T", stmt)
	switch stmtType := stmt.(type) {
	case grammar.Token:
		return makeTemplateStr(offset, nodeType, fmt.Sprintf("type [%v] lexeme [%v] literal [%v]", stmtType.TokenType, stmtType.Lexeme, stmtType.Literal))
	case grammar.VariableDeclarationStatement:
		initExpr := printer.printNode(offset+1, stmtType.Initializer)
		token := printer.printNode(offset+1, stmtType.Name)
		return makeTemplateStr(offset, nodeType, token, initExpr)
	case grammar.PrintStatement:
		value := printer.printNode(offset+1, stmtType.Value)
		return makeTemplateStr(offset, nodeType, value)
	case grammar.BlockScopeStatement:
		stmts := printer.printNode(offset+1, stmtType.Statements)
		return makeTemplateStr(offset, nodeType, stmts)
	case grammar.WhileLoopStatement:
		expr := printer.printNode(offset+1, stmtType.Condition)
		body := printer.printNode(offset+1, stmtType.Body)
		return makeTemplateStr(offset, nodeType, expr, body)
	case grammar.ReturnStatement:
		expr := printer.printNode(offset+1, stmtType.Expression)
		keyword := printer.printNode(offset+1, stmtType.Keyword)
		return makeTemplateStr(offset, nodeType, expr, keyword)
	case grammar.ConditionalStatement:
		condition := printer.printNode(offset+1, stmtType.Condition)
		thenBranch := printer.printNode(offset+1, stmtType.ThenBranch)
		elseBranch := printer.printNode(offset+1, stmtType.ElseBranch)
		return makeTemplateStr(offset, nodeType, condition, thenBranch, elseBranch)
	case grammar.ClassDeclarationStatement:
		className := printer.printNode(offset+1, stmtType.Name)
		var builder strings.Builder
		for _, method := range stmtType.Methods {
			builder.WriteString(printer.printNode(offset+1, method))
		}
		methods := builder.String()
		return makeTemplateStr(offset, nodeType, className, methods)
	case []grammar.Statement:
		var builder strings.Builder
		for _, statement := range stmtType {
			builder.WriteString(printer.printNode(offset+1, statement))
		}
		return builder.String()
	case grammar.ExpressionStatement:
		expr := printer.printNode(offset+1, stmtType.Expression)
		return makeTemplateStr(offset, nodeType, expr)
	case grammar.UnaryExpression:
		token := printer.printNode(offset+1, stmtType.Operator)
		rightExpr := printer.printNode(offset+1, stmtType.Right)
		return makeTemplateStr(offset, nodeType, token, rightExpr)
	case grammar.BinaryExpression:
		leftExpr := printer.printNode(offset+1, stmtType.Left)
		operator := printer.printNode(offset+1, stmtType.Operator)
		rightExpr := printer.printNode(offset+1, stmtType.Right)
		return makeTemplateStr(offset, nodeType, leftExpr, operator, rightExpr)
	case grammar.LiteralExpression:
		literal := fmt.Sprintf("literal [%v]", stmtType.Literal)
		return makeTemplateStr(offset, nodeType, literal)
	case grammar.VariableDeclaration:
		token := printer.printNode(offset+1, stmtType.Name)
		return makeTemplateStr(offset, nodeType, token)
	case grammar.FunctionDeclarationStatement:
		token := printer.printNode(offset+1, stmtType.Name)
		body := printer.printNode(offset+1, stmtType.Body)
		return makeTemplateStr(offset, nodeType, token, body)
	case grammar.LogicExpression:
		leftExpr := printer.printNode(offset+1, stmtType.Left)
		operator := printer.printNode(offset+1, stmtType.Operator)
		rightExpr := printer.printNode(offset+1, stmtType.Right)
		return makeTemplateStr(offset, nodeType, leftExpr, operator, rightExpr)
	case grammar.GroupingExpression:
		expr := printer.printNode(offset+1, stmtType.Expression)
		return makeTemplateStr(offset, nodeType, expr)
	case grammar.AssignmentExpression:
		token := printer.printNode(offset+1, stmtType.Name)
		expr := printer.printNode(offset+1, stmtType.Value)
		return makeTemplateStr(offset, nodeType, token, expr)
	case grammar.CallExpression:
		callee := printer.printNode(offset+1, stmtType.Callee)
		expr := printer.printNode(offset+1, stmtType.Paren)
		args := printer.printNode(offset+1, stmtType.Arguments)
		return makeTemplateStr(offset, nodeType, callee, expr, args)
	case grammar.PropertyAccessExpression:
		object := printer.printNode(offset+1, stmtType.Object)
		return makeTemplateStr(offset, nodeType, object)
	case grammar.PropertyAssignmentExpression:
		object := printer.printNode(offset+1, stmtType.Object)
		value := printer.printNode(offset+1, stmtType.Value)
		return makeTemplateStr(offset, nodeType, object, value)
	case grammar.SelfReferenceExpression:
		keyword := printer.printNode(offset+1, stmtType.Keyword)
		return makeTemplateStr(offset, nodeType, keyword)
	case grammar.BaseClassCallExpression:
		keyword := printer.printNode(offset+1, stmtType.Keyword)
		method := printer.printNode(offset+1, stmtType.Method)
		return makeTemplateStr(offset+1, nodeType, keyword, method)
	default:
		return nodeType
	}
}

func makeTemplateStr(offset int, args ...string) string {
	var builder strings.Builder
	for index, arg := range args {
		if index == 0 {
			builder.WriteString(fmt.Sprintf("[%s] => {\n", arg))
		} else {
			builder.WriteString(offsetTemplateStr("", offset))
			builder.WriteString(fmt.Sprintf("%s\n", arg))
		}
	}
	builder.WriteString(offsetTemplateStr("", offset))
	builder.WriteString("}")
	return builder.String()
}

func offsetTemplateStr(str string, offset int) string {
	return fmt.Sprintf("%*s", offset, str)
}

func (printer *TokenPrinter) Print(tokens []grammar.Token) {
	fmt.Println("Tokens generated from source:")
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for index, token := range tokens {
		fmt.Fprintln(writer, printer.printToken(index, token))
	}
	writer.Flush()
	fmt.Println()
}

func (printer *TokenPrinter) printToken(index int, token grammar.Token) string {
	return fmt.Sprintf("[%v]\t %T => type [%v]\t lexeme [%v]\t literal [%v]", index, token, token.TokenType, token.Lexeme, token.Literal)
}
