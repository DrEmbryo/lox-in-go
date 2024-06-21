package grammar

type Expression any

type BinaryExpression struct {
	Left     Expression
	Operator Token
	Right    Expression
}

type UnaryExpression struct {
	Operator Token
	Right    Expression
}

type LiteralExpression struct {
	Literal any
}

type GroupingExpression struct {
	Expression Expression
}