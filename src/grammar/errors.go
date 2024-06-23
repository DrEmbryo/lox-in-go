package grammar

type LoxError interface {
	Error() string
	Print()
}