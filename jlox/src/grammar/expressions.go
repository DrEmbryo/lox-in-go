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

type PropertyAccessExpression struct {
	Object Expression
	Name   Token
}

type PropertyAssignmentExpression struct {
	Object Expression
	Value  Expression
	Name   Token
}

type SelfReferenceExpression struct {
	Keyword Token
}

type BaseClassCallExpression struct {
	Keyword Token
	Method  Token
}
