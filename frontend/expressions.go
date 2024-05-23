package Lox

type Expression interface{}

type BinaryExpression struct {
	left     Expression
	operator Token
	right    Expression
}

type UnaryExpression struct {
	operator Token
	right    Expression
}

type LiteralExpression struct {
	literal interface{}
}

type GroupingExpression struct {
	expression Expression
}