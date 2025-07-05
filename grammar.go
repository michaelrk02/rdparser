package rdparser

import (
	"context"
	"fmt"

	"github.com/shivamMg/rd"
)

var ErrCompile = fmt.Errorf("compile error")

type Grammar interface {
	BuildParseTree(ctx context.Context, b *Builder) error
}

func Compile(tokens []rd.Token, g Grammar) (*Tree, error) {
	b := &Builder{Builder: rd.NewBuilder(tokens)}

	err := g.BuildParseTree(context.Background(), b)
	if err != nil {
		return nil, err
	}

	return &Tree{Tree: b.ParseTree()}, nil
}

func NewSyntaxError(ctx context.Context, b *Builder) error {
	return NewError(ErrCompile, fmt.Sprintf("invalid syntax near token `%s`", b.Last()))
}
