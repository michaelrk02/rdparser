package main

import (
	"context"
	"flag"
	"fmt"
	"math"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula"
)

func main() {
	var expr string
	var epsilon float64

	flag.StringVar(&expr, "expr", "", "expression")
	flag.Float64Var(&epsilon, "epsilon", 0.0, "use this epsilon (error-tolerance) value")
	flag.Parse()

	if expr == "" {
		flag.PrintDefaults()
		return
	}

	varDict := formula.VariableDict{
		"pi":  math.Pi,
		"e":   math.E,
		"nan": math.NaN(),
		"inf": math.Inf(1),
	}

	lexer := formula.NewLexer()
	grammar := formula.NewGrammar()

	stdlib := formula.NewStdLibrary()
	parser := formula.NewParser(stdlib, epsilon, varDict)

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
		panic(err)
	}

	n := rslt.(float64)
	fmt.Printf("%v\n", n)
}
