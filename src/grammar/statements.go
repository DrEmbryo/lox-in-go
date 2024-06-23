package grammar

type Statement any

type ExpressionStatement struct {
	Expression any
}

type PrintStatment struct {
	Value Expression
}
