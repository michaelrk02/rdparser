package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula"
)

func main() {
	var expr string
	var test string
	var epsilon float64

	flag.StringVar(&expr, "expr", "", "expression")
	flag.StringVar(&test, "test", "", "load test case from file")
	flag.Float64Var(&epsilon, "epsilon", 0.0, "use this epsilon value")
	flag.Parse()

	if expr == "" && test == "" {
		flag.PrintDefaults()
		return
	}

	tokenizer := formula.NewTokenizer()
	grammar := formula.NewGrammar()

	stdlib := formula.NewStdLibrary()
	parser := formula.NewParser(stdlib.Library, epsilon)

	if expr != "" {
		tokens, err := tokenizer.Tokenize(expr)
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
	} else if test != "" {
		loadTestCase(test, tokenizer, grammar, parser)
	}
}

func loadTestCase(file string, tokenizer *formula.Tokenizer, grammar *formula.Grammar, parser *formula.Parser) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	lineNum := 1
	for {
		row, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				panic(err)
			}
		}

		if lineNum == 1 {
			lineNum++
			continue
		}

		expr := row[0]
		expected, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			expected = math.NaN()
		}

		tokens, err := tokenizer.Tokenize(expr)
		if err != nil {
			panic(err)
		}

		tree, err := rdparser.Compile(tokens, grammar)
		if err != nil {
			panic(err)
		}

		rslt, err := parser.Parse(context.Background(), tree)
		if err != nil {
			rslt = math.NaN()
		}

		actual := rslt.(float64)

		ok := parser.IsEqu(expected, actual)

		msg := "FAIL"
		if ok {
			msg = "OK"
		}

		fmt.Printf("[%s] [line %d] %s [expected:%v actual:%v]\n", msg, lineNum, expr, expected, actual)

		lineNum++
	}
}
