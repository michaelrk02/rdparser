package rdparser

import (
	"context"

	"github.com/shivamMg/rd"
)

type Builder struct {
	*rd.Builder

	last rd.Token
}

func (b *Builder) Enter(ctx *context.Context, sym NonTerminal) *Builder {
	b.Builder.Enter(sym)
	if ctx != nil {
		*ctx = Trace(*ctx, sym)
	}
	return b
}

func (b *Builder) Match(token rd.Token) bool {
	if b.Builder.Match(token) {
		b.last = token
		return true
	}
	return false
}

func (b *Builder) Add(token rd.Token) {
	b.Builder.Add(token)
	b.last = token
}

func (b *Builder) Last() rd.Token {
	return b.last
}
