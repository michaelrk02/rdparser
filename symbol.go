package rdparser

import (
	"fmt"
	"strings"
)

type Terminal string

func (t Terminal) String() string {
	return string(t)
}

func (t Terminal) Hex() string {
	s := t.String()
	hexArr := make([]string, len(s))
	for i, c := range s {
		hexArr[i] = fmt.Sprintf("\\x%02x", c)
	}
	return strings.Join(hexArr, "")
}

type NonTerminal string

func (t NonTerminal) String() string {
	return string(t)
}

func IsTerminal(v interface{}) bool {
	if _, ok := v.(Terminal); ok {
		return true
	}
	return false
}

func IsTerminalOf(v interface{}, sym Terminal) bool {
	if test, ok := v.(Terminal); ok && test == sym {
		return true
	}
	return false
}

func IsNonTerminal(v interface{}) bool {
	if _, ok := v.(NonTerminal); ok {
		return true
	}
	return false
}

func IsNonTerminalOf(v interface{}, sym NonTerminal) bool {
	if test, ok := v.(NonTerminal); ok && test == sym {
		return true
	}
	return false
}
