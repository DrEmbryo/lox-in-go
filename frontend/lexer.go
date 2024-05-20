package Lox

const (
	// single char tokens
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE 
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// multi char tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords 
	AND
	CLASS
	ELSE
	FALSE
	FUNC
	FOR
	IF
	NULL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

type Token struct {
	tokenType uint8
	lexeme string
	literal string
}

func (token Token) Create (tokenType uint8, lexeme string, literal string) (Token) {
	token.tokenType = tokenType
	token.lexeme = lexeme
	token.literal = literal
	return token
}

type Lexer struct {
	start uint32
	current uint32
	line uint32
}

func (lexer *Lexer) Next (source []rune) rune {
	nextChar := source[lexer.current]
	lexer.current ++
	return nextChar
}


func (lexer Lexer) Tokenize(source string) ([]Token, []Error) {
	lexErrors := make([]Error,0)
	tokens := make([]Token, 10)
	lexer.start = 1
	lexer.current = 0
	lexer.line = 0

	for  {
		char := lexer.Next([]rune(source))
		if lexer.current >= uint32(len(source)) {
			break
		}

		switch char {
		case '(':
			tokens = append(tokens, Token{}.Create(LEFT_PAREN, "", ""))
		case ')':
			tokens = append(tokens, Token{}.Create(RIGHT_PAREN, "", ""))
		case '{':
			tokens = append(tokens, Token{}.Create(LEFT_BRACE, "", ""))
		case '}':
			tokens = append(tokens, Token{}.Create(RIGHT_BRACE, "", ""))
		case ',':
			tokens = append(tokens, Token{}.Create(COMMA, "", ""))
		case '.':
			tokens = append(tokens, Token{}.Create(DOT, "", ""))
		case '-':
			tokens = append(tokens, Token{}.Create(MINUS, "", ""))
		case '+':
			tokens = append(tokens, Token{}.Create(PLUS, "", ""))
		case ';':
			tokens = append(tokens, Token{}.Create(SEMICOLON, "", ""))
		case '*':
			tokens = append(tokens, Token{}.Create(STAR, "", ""))
		default:
			lexErrors = append(lexErrors, Error{ line: lexer.line, position: lexer.current, message: "Unidentified token"})
		}

	}
	
	tokens = append(tokens, Token{ 
		tokenType: EOF,
		lexeme: "",
		literal: "",
	})
	return tokens, lexErrors
}
