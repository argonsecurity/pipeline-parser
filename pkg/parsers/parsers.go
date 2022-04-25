package parsers

type Parser[T any] interface {
	Parse(data []byte) (*T, error)
}
