package grammar

type Statement any

type ExpressionStatement struct {
	Expression Expression
}

type PrintStatement struct {
	Value Expression
}

type VariableDeclarationStatement struct {
	Name        Token
	Initializer Expression
}

type BlockScopeStatement struct {
	Statements []Statement
}

type ConditionalStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

type WhileLoopStatement struct {
	Expression Expression
	Body       Statement
}
