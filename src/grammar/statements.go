package grammar

type Statement any

type ExpressionStatement struct {
	Expression Expression
}

type PrintStatment struct {
	Value Expression
}

type VariableDeclarationStatment struct {
	Name        Token
	Initializer Expression
}