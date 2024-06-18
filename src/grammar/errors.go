package grammar

type LoxError interface {
	Error() string
	Print()
}

// type ParserError struct {
// 	Position uint32
// 	Token Token
// 	Message	string
// }

// func (e ParserError) Print() {
// 	if (e.Token.TokenType == EOF) {
// 		fmt.Printf("[%d]: EOF at %s, %s", e.Position, e.Token.Lexeme, e.Message)
// 	} else {
// 		fmt.Printf("[%d]: Token %d: failed at %s with value %s: %s", e.Position, e.Token.TokenType, e.Token.Lexeme, e.Token.Literal, e.Message)
// 	}
// }

// func (e ParserError) Error() string {
// 	return fmt.Sprintf("[%d]: Token %d: failed at %s with value %s: %s", e.Position, e.Token.TokenType, e.Token.Lexeme, e.Token.Literal, e.Message)
// }

// type RuntimeError struct {
// 	Token Token
// 	Message string
// }

// func (e RuntimeError) Print() {
// 	fmt.Printf("[%v]: Runtime error: %s", e.Token, e.Message)
// }

// func (e RuntimeError) Error() string {
// 	return fmt.Sprintf("[%v]: Runtime error: %s", e.Token, e.Message)
// }