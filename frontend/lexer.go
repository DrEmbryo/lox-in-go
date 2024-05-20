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
	lexeme    string
	literal   string
}

func (token Token) Create(tokenType uint8, lexeme string, literal string) Token {
	token.tokenType = tokenType
	token.lexeme = lexeme
	token.literal = literal
	return token
}

type Lexer struct {
	start   uint32
	current uint32
	line    uint32
	source  []rune
}

func (lexer *Lexer) init(source []rune, start uint32, current uint32, line uint32){
	lexer.start = start
	lexer.current = current
	lexer.line = line
	lexer.source = source
}

func (lexer *Lexer) next() rune {
	if lexer.current < uint32(len(lexer.source)) {
		nextChar := lexer.source[lexer.current]
		lexer.current++
		return nextChar
	}
	return 0
}

func (lexer Lexer) lookahead(offset int) rune {
	if lexer.current < uint32(len(lexer.source) + offset) {
		return lexer.source[lexer.current]
	}
	return 0
}

func (lexer Lexer) Tokenize(source string) ([]Token, []Error) {
	tokens := make([]Token, 10)
	lexer.init([]rune(source), 1, 0, 0)
	lexErrors := make([]Error, 0)

	for {
		char := lexer.next()
		if lexer.current >= uint32(len(source)) {
			break
		}
		// fmt.Println(string(char))
		var tokenType uint8
		switch char {
		case '(':
			tokenType = LEFT_PAREN
		case ')':
			tokenType = RIGHT_PAREN
		case '{':
			tokenType = LEFT_BRACE
		case '}':
			tokenType = RIGHT_BRACE
		case ',':
			tokenType = COMMA
		case '.':
			tokenType = DOT
		case '-':
			tokenType = MINUS
		case '+':
			tokenType = PLUS
		case ';':
			tokenType = SEMICOLON
		case '*':
			lookahead := lexer.lookahead(0)
			if lookahead == '/' {
				// multiline comment skip
				for {	
					next := lexer.next()
					lookup := lexer.lookahead(1)
					if next == 0 || lookup == 0 {
						break
					} 
					if next == '/' && lookup == '*' {
						lexer.next()
						break
					}
				}
			} else {
				tokenType = STAR
			}
		case '/':
			lookahead := lexer.lookahead(0)
			if lookahead == '/' {
				// comment line skip
				for {
					if lexer.next() == '\n' || lexer.next() == 0{
						break
					}
				}
			} else {
				tokenType = SLASH
			}
		case '=':
			
			if lexer.next() == '=' {
				tokenType = EQUAL_EQUAL
			} else {
				tokenType = EQUAL
			}
		case '!':
			if lexer.next() == '=' {
				tokenType = BANG_EQUAL
			} else {
				tokenType = BANG
			}
		case '<':
			if lexer.next() == '=' {
				tokenType = LESS_EQUAL
			} else {
				tokenType = EQUAL
			}
		case '>':
			if lexer.next() == '=' {
				tokenType = GREATER_EQUAL
			} else {
				tokenType = EQUAL
			}
		case ' ':
		case '\r':
		case '\t':
		case '\n': 
		lexer.line++
		default:
			lexErrors = append(lexErrors, Error{line: lexer.line, position: lexer.current, message: "Unidentified token"})
			tokenType = 0
		}

		tokens = append(tokens, Token{}.Create(tokenType, "", ""))
	}

	tokens = append(tokens, Token{
		tokenType: EOF,
		lexeme:    "",
		literal:   "",
	})
	return tokens, lexErrors
}
