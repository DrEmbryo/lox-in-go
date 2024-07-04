package utils

import (
	"fmt"
	"os"
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
		printer.leftPad = 0
		fmt.Println(printer.printNode(stmt))
	}
	fmt.Println("")
}

func (printer *AstPrinter) printNode(stmt grammar.Statement) string {
	printer.leftPad++
	switch stmtType := stmt.(type) {
	case grammar.VariableDeclarationStatement:
		return fmt.Sprintf("%T =>\n token [%v]\n init expression [%v]\n", stmt, printer.tokenPrinter.printToken(stmtType.Name), printer.printNode(stmtType.Initializer))
	case grammar.PrintStatement:
		return fmt.Sprintf("%T =>\n value [%v]\n", stmtType, printer.printNode(stmtType.Value))
	case grammar.BlockScopeStatement:
		return fmt.Sprintf("%T =>\n ", stmtType)
	case grammar.ConditionalStatement:
		return fmt.Sprintf("%T =>\n condition [%v]\n then branch [%v]\n else branch [%v]\n", stmtType, printer.printNode(stmtType.Condition), printer.printNode(stmtType.ThenBranch), printer.printNode(stmtType.ElseBranch))
	case grammar.ExpressionStatement:
		return fmt.Sprintf("%T =>\n expression [%v]\n", stmtType, printer.printNode(stmtType.Expression))
	case grammar.UnaryExpression:
		return fmt.Sprintf("%T =>\n operator [%v]\n right hand expression [%v]\n", stmtType, printer.tokenPrinter.printToken(stmtType.Operator), printer.printNode(stmtType.Right))
	case grammar.BinaryExpression:
		return fmt.Sprintf("%T =>\n left hand expression [%v]\n operator [%v]\n right hand expression [%v]", stmtType, printer.printNode(stmtType.Left), printer.tokenPrinter.printToken(stmtType.Operator), printer.printNode(stmtType.Right))
	case grammar.LiteralExpression:
		return fmt.Sprintf("%T => literal [%v]", stmtType, stmtType.Literal)
	case grammar.VariableDeclaration:
		return fmt.Sprintf("%T => token [%v]", stmtType, printer.tokenPrinter.printToken(stmtType.Name))
	case grammar.GroupingExpression:
		return fmt.Sprintf("%T =>\n expression [%v]\n", stmtType, printer.printNode(stmtType.Expression))
	default:
		return fmt.Sprintf("%T", stmtType)
	}
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
	return fmt.Sprintf("type [%v]\t lexeme [%v]\t literal [%v]", token.TokenType, token.Lexeme, token.Literal)
}
