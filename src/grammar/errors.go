package grammar

type LoxError interface {
	Error() string
	Print()
}

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