package formula

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula/logic"
)

const (
	Epsilon = 0.0

	TestcaseFile = "testcases/example.csv"
)

func TestFormula(t *testing.T) {
	varDict := VariableDict{}

	lexer := NewLexer()
	grammar := NewGrammar()

	lib := NewTestLib()
	parser := NewParser(lib, Epsilon, varDict)

	f, err := os.Open(TestcaseFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	for lineNum := 1; ; lineNum++ {
		row, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				panic(err)
			}
		}

		expr := row[0]
		expected, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			t.Error(err)
			continue
		}

		tokens, err := lexer.Lex(expr)
		if err != nil {
			panic(err)
		}

		tree, err := rdparser.Compile(tokens, grammar)
		if err != nil {
			panic(err)
		}

		rslt, err := parser.Parse(context.Background(), tree)
		if err != nil {
			t.Error(err)
			continue
		}

		actual := rslt.(float64)

		ok := logic.Equ(expected, actual, Epsilon)

		msg := "FAIL"
		if ok {
			msg = "OK"
		} else {
			t.Fail()
		}

		t.Logf("[%-4s : %-4d] %s [expected:%v actual:%v]", msg, lineNum, expr, expected, actual)
	}
}

type TestLib struct {
	*StdLibrary

	ref map[string]Function
}

func NewTestLib() Library {
	lib := &TestLib{
		StdLibrary: NewStdLibrary(),
		ref:        make(map[string]Function),
	}
	return lib
}

func (lib *TestLib) Resolve(funcName string) (Function, bool) {
	if fn, ok := lib.ref[funcName]; ok {
		return fn, true
	}
	return lib.StdLibrary.Resolve(funcName)
}
