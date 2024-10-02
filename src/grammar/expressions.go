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

type VariableDeclaration struct {
	Name Token
}

type AssignmentExpression struct {
	Name  Token
	Value Expression
}

type LogicExpression struct {
	Left     Expression
	Operator Token
	Right    Expression
}

type CallExpression struct {
	Callee    Expression
	Paren     Token
	Arguments []Expression
}
