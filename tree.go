package rdparser

import (
	"fmt"
	"strings"

	"github.com/shivamMg/rd"
)

type Tree struct {
	*rd.Tree
}

func (t *Tree) Len() int {
	return len(t.Subtrees)
}

func (t *Tree) Has(n int) bool {
	return t.Len() == n
}

func (t *Tree) AsTerminal() Terminal {
	if sym, ok := t.Symbol.(Terminal); ok {
		return sym
	}
	panic(NewRuntimeError("not a terminal symbol"))
}

func (t *Tree) AsNonTerminal() NonTerminal {
	if sym, ok := t.Symbol.(NonTerminal); ok {
		return sym
	}
	panic(NewRuntimeError("not a non-terminal symbol"))
}

func (t *Tree) At(index int) *Tree {
	return &Tree{Tree: t.Subtrees[index]}
}

func (t *Tree) IsTerminal() bool {
	return IsTerminal(t.Symbol)
}

func (t *Tree) IsTerminalOf(sym Terminal) bool {
	return IsTerminalOf(t.Symbol, sym)
}

func (t *Tree) IsNonTerminal() bool {
	return IsNonTerminal(t.Symbol)
}

func (t *Tree) IsNonTerminalOf(sym NonTerminal) bool {
	return IsNonTerminalOf(t.Symbol, sym)
}

func (t *Tree) AssertTerminal() *Tree {
	if !t.IsTerminal() {
		panic(NewRuntimeError("expecting terminal symbol"))
	}
	return t
}

func (t *Tree) AssertTerminalOf(sym Terminal) *Tree {
	if !t.IsTerminalOf(sym) {
		panic(NewRuntimeError(fmt.Sprintf("invalid terminal symbol, expecting `%s`", sym)))
	}
	return t
}

func (t *Tree) AssertNonTerminal() *Tree {
	if !t.IsNonTerminal() {
		panic(NewRuntimeError("expecting non-terminal symbol"))
	}
	return t
}

func (t *Tree) AssertNonTerminalOf(sym NonTerminal) *Tree {
	if !t.IsNonTerminalOf(sym) {
		panic(NewRuntimeError(fmt.Sprintf("invalid non-terminal symbol, expecting `%s`", sym)))
	}
	return t
}

func (t *Tree) Traverse(sep string) string {
	return strings.Join(t.Walk(), sep)
}

func (t *Tree) Walk() []string {
	if t.Len() == 0 {
		return []string{string(t.AsTerminal())}
	}

	children := []string{}
	for _, sub := range t.Subtrees {
		subt := &Tree{Tree: sub}
		children = append(children, subt.Walk()...)
	}

	return children
}
