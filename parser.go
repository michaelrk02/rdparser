package rdparser

import (
	"context"
	"fmt"
)

var ErrParse = fmt.Errorf("parse error")

type Parser interface {
	Parse(ctx context.Context, t *Tree) (interface{}, error)
}

func NewParseError(ctx context.Context, msg string) error {
	return NewError(ErrParse, fmt.Sprintf("%s (stacktrace = %s)", msg, GetStackTrace(ctx)))
}
