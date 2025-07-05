package rdparser

import (
	"fmt"

	"github.com/shivamMg/rd"
)

var ErrLexical = fmt.Errorf("lexical error")

type Lexer interface {
	Lex(input string) ([]rd.Token, error)
}

func NewLexicalError(msg string) error {
	return NewError(ErrLexical, msg)
}
