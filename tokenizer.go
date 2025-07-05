package rdparser

import (
	"fmt"

	"github.com/shivamMg/rd"
)

var ErrTokenize = fmt.Errorf("tokenize error")

type Tokenizer interface {
	Tokenize(input string) ([]rd.Token, error)
}
