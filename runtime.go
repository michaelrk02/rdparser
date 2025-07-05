package rdparser

import (
	"context"
	"fmt"
	"strings"
)

var (
	ErrRuntime = fmt.Errorf("runtime error")

	keyStackTrace = contextKey("StackTrace")
)

func NewRuntimeError(msg string) error {
	return NewError(ErrRuntime, msg)
}

func Trace(ctx context.Context, symbol NonTerminal) context.Context {
	newStackTrace := StackTrace{
		Path: append(GetStackTrace(ctx).Path, StackTraceElement{
			Symbol: symbol,
			Data:   "",
		}),
	}
	return context.WithValue(ctx, keyStackTrace, newStackTrace)
}

type StackTrace struct {
	Path []StackTraceElement
}

type StackTraceElement struct {
	Symbol NonTerminal
	Data   interface{}
}

func GetStackTrace(ctx context.Context) StackTrace {
	if v := ctx.Value(keyStackTrace); v != nil {
		if st, ok := v.(StackTrace); ok {
			return st
		}
	}
	return StackTrace{Path: []StackTraceElement{}}
}

func (st StackTrace) Lookup(sym NonTerminal) (bool, int) {
	for i := range st.Path {
		if st.Path[len(st.Path)-1-i].Symbol == sym {
			return true, i
		}
	}
	return false, 0
}

func (st StackTrace) String() string {
	elems := make([]string, len(st.Path))
	for i := range st.Path {
		elems[i] = st.Path[i].Symbol.String()
	}
	return strings.Join(elems, " > ")
}
