package utils

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/DrEmbryo/lox/src/grammar"
)

type AstPrinter struct {
	leftPad      int
	tokenPrinter TokenPrinter
}

type TokenPrinter struct {
}

func (printer *AstPrinter) Print(stmts []grammar.Statement) {
	fmt.Println("Ast generated from tokens:")
	for _, stmt := range stmts {
		printer.leftPad = -1
		fmt.Println(printer.printNode(stmt))
	}
	fmt.Println("")
}

func (printer *AstPrinter) printNode(stmt grammar.Statement) string {
	nodeType := fmt.Sprintf("%T", stmt)
	switch stmtType := stmt.(type) {
	case grammar.VariableDeclarationStatement:
		token := printer.tokenPrinter.printToken(stmtType.Name)
		initExpr := printer.printNode(stmtType.Initializer)
		return makeTemplateStr(nodeType, token, initExpr)
	case grammar.PrintStatement:
		value := printer.printNode(stmtType.Value)
		return makeTemplateStr(nodeType, value)
	case grammar.BlockScopeStatement:
		stmts := printer.printNode(stmtType.Statements)
		return makeTemplateStr(nodeType, stmts)
	case grammar.ConditionalStatement:
		condition := printer.printNode(stmtType.Condition)
		thenBranch := printer.printNode(stmtType.ThenBranch)
		elseBranch := printer.printNode(stmtType.ElseBranch)
		return makeTemplateStr(nodeType, condition, thenBranch, elseBranch)
	case grammar.ExpressionStatement:
		expr := printer.printNode(stmtType.Expression)
		return makeTemplateStr(nodeType, expr)
	case grammar.UnaryExpression:
		token := printer.tokenPrinter.printToken(stmtType.Operator)
		rightExpr := printer.printNode(stmtType.Right)
		return makeTemplateStr(nodeType, token, rightExpr)
	case grammar.BinaryExpression:
		leftExpr := printer.printNode(stmtType.Left)
		operator := printer.tokenPrinter.printToken(stmtType.Operator)
		rightExpr := printer.printNode(stmtType.Right)
		return makeTemplateStr(nodeType, leftExpr, operator, rightExpr)
	case grammar.LiteralExpression:
		literal := fmt.Sprintf("literal [%v]", stmtType.Literal)
		return makeTemplateStr(nodeType, literal)
	case grammar.VariableDeclaration:
		token := printer.tokenPrinter.printToken(stmtType.Name)
		return makeTemplateStr(nodeType, token)
	case grammar.GroupingExpression:
		expr := printer.printNode(stmtType.Expression)
		return makeTemplateStr(nodeType, expr)
	case grammar.AssignmentExpression:
		token := printer.tokenPrinter.printToken(stmtType.Name)
		expr := printer.printNode(stmtType.Value)
		return makeTemplateStr(nodeType, token, expr)
	default:
		return nodeType
	}
}

func makeTemplateStr(args ...string) string {
	var builder strings.Builder
	for index, arg := range args {
		if index == 0 {
			builder.WriteString(fmt.Sprintf("[%s] => { ", arg))
		} else {
			builder.WriteString(fmt.Sprintf("%s ", arg))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

func (printer *TokenPrinter) Print(tokens []grammar.Token) {
	fmt.Println("Tokens generated from source:")
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, token := range tokens {
		fmt.Fprintln(writer, printer.printToken(token))
	}
	writer.Flush()
	fmt.Println()
}

func (printer *TokenPrinter) printToken(token grammar.Token) string {
	return fmt.Sprintf("%T => type [%v]\t lexeme [%v]\t literal [%v]", token, token.TokenType, token.Lexeme, token.Literal)
}
