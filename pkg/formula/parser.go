package formula

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula/symbol"
	"github.com/michaelrk02/rdparser/pkg/formula/token"
)

type Parser struct {
	lib     *Library
	epsilon float64

	Var      map[string]float64
	varRegex *regexp.Regexp
}

func NewParser(lib *Library, epsilon float64) *Parser {
	return &Parser{
		lib:      lib,
		epsilon:  epsilon,
		Var:      make(map[string]float64),
		varRegex: regexp.MustCompile(`^\[([a-zA-Z0-9-:]+)\]$`),
	}
}

func (p *Parser) Parse(ctx context.Context, t *rdparser.Tree) (rslt interface{}, err error) {
	defer rdparser.Catch(rdparser.ErrParse, &err)

	rslt = p.Expr(ctx, t)
	return
}

func (p *Parser) Expr(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.Expr)

	term := p.Term(ctx, t.At(0).AssertNonTerminalOf(symbol.Term))

	if t.At(1).AssertNonTerminalOf(symbol.Exprx).Has(2) {
		op := t.At(1).At(0).AsTerminal()
		switch op {
		case token.Add:
			term = term + p.Expr(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.Expr))
		case token.Sub:
			term = term - p.Expr(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.Expr))
		}
	}

	return term
}

func (p *Parser) Term(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.Term)

	factor := p.Factor(ctx, t.At(0).AssertNonTerminalOf(symbol.Factor))

	if t.At(1).AssertNonTerminalOf(symbol.Termx).Has(2) {
		op := t.At(1).At(0).AssertTerminal().AsTerminal()
		switch op {
		case token.Mul:
			factor = factor * p.Term(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.Term))
		case token.Div:
			factor = factor / p.Term(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.Term))
		case token.Mod:
			factor = float64(int(factor) % int(p.Term(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.Term))))
		}
	}

	return factor
}

func (p *Parser) Factor(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.Factor)

	if t.At(0).IsTerminalOf(token.LParen) && t.At(1).IsNonTerminalOf(symbol.Factorx) {
		if t.At(1).At(0).IsNonTerminalOf(symbol.BoolCond) {
			return p.BoolCond(ctx, t.At(1).At(0))
		}

		if t.At(1).At(0).IsNonTerminalOf(symbol.Expr) {
			return p.Expr(ctx, t.At(1).At(0))
		}
	}

	if t.At(0).IsTerminalOf(token.Minus) && t.At(1).IsNonTerminalOf(symbol.Factor) {
		return -p.Factor(ctx, t.At(1))
	}

	if t.At(0).IsNonTerminalOf(symbol.Variable) {
		return p.Variable(ctx, t.At(0))
	}

	if t.At(0).IsNonTerminalOf(symbol.Number) {
		return p.Number(ctx, t.At(0))
	}

	if t.At(0).IsNonTerminalOf(symbol.FuncCall) {
		return p.FuncCall(ctx, t.At(0))
	}

	panic(rdparser.NewParseError(ctx, "invalid expression"))
}

func (p *Parser) FuncCall(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.FuncCall)

	funcName := t.At(0).AssertNonTerminalOf(symbol.FuncName).At(0).AsTerminal().String()

	t.At(1).AssertTerminalOf(token.LParen)
	t.At(2).AssertNonTerminalOf(symbol.FuncArg)
	t.At(3).AssertTerminalOf(token.RParen)

	funcArgs := p.FuncArg(ctx, t.At(2))

	if callback, ok := p.lib.Data[funcName]; ok {
		return callback(ctx, funcArgs)
	}

	panic(rdparser.NewParseError(ctx, fmt.Sprintf("unrecognized function `%s`", funcName)))
}

func (p *Parser) FuncArg(ctx context.Context, t *rdparser.Tree) []float64 {
	ctx = rdparser.Trace(ctx, symbol.FuncArg)

	if t.Has(2) && t.At(0).IsNonTerminalOf(symbol.Expr) && t.At(1).IsNonTerminalOf(symbol.FuncArgx) {
		args := []float64{p.Expr(ctx, t.At(0))}

		if t.At(1).Has(2) && t.At(1).At(0).IsTerminalOf(token.Comma) && t.At(1).At(1).IsNonTerminalOf(symbol.FuncArg) {
			args = append(args, p.FuncArg(ctx, t.At(1).At(1))...)
		}

		return args
	}

	return []float64{}
}

func (p *Parser) BoolCond(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.BoolCond)

	boolExpr := p.BoolExpr(ctx, t.At(0).AssertNonTerminalOf(symbol.BoolExpr))
	t.At(1).AssertTerminalOf(token.Question)
	t.At(2).AssertNonTerminalOf(symbol.Expr)
	t.At(3).AssertTerminalOf(token.Colon)
	t.At(4).AssertNonTerminalOf(symbol.Expr)

	if boolExpr {
		return p.Expr(ctx, t.At(2))
	} else {
		return p.Expr(ctx, t.At(4))
	}
}

func (p *Parser) BoolExpr(ctx context.Context, t *rdparser.Tree) bool {
	ctx = rdparser.Trace(ctx, symbol.BoolExpr)

	boolTerm := p.BoolTerm(ctx, t.At(0).AssertNonTerminalOf(symbol.BoolTerm))

	if t.At(1).AssertNonTerminalOf(symbol.BoolExprx).Has(2) {
		t.At(1).At(0).AssertNonTerminalOf(symbol.LogicOr)
		return boolTerm || p.BoolExpr(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.BoolExpr))
	}

	return boolTerm
}

func (p *Parser) BoolTerm(ctx context.Context, t *rdparser.Tree) bool {
	ctx = rdparser.Trace(ctx, symbol.BoolTerm)

	boolFactor := p.BoolFactor(ctx, t.At(0).AssertNonTerminalOf(symbol.BoolFactor))

	if t.At(1).AssertNonTerminalOf(symbol.BoolTermx).Has(2) {
		t.At(1).At(0).AssertNonTerminalOf(symbol.LogicAnd)
		return boolFactor && p.BoolTerm(ctx, t.At(1).At(1).AssertNonTerminalOf(symbol.BoolTerm))
	}

	return boolFactor
}

func (p *Parser) BoolFactor(ctx context.Context, t *rdparser.Tree) bool {
	ctx = rdparser.Trace(ctx, symbol.BoolFactor)

	if t.Has(2) && t.At(0).IsNonTerminalOf(symbol.LogicNot) {
		return !p.BoolFactor(ctx, t.At(1).AssertNonTerminalOf(symbol.BoolFactor))
	}

	if t.Has(3) && t.At(0).IsTerminalOf(token.LParen) {
		t.At(2).AssertTerminalOf(token.RParen)

		return p.BoolExpr(ctx, t.At(1).AssertNonTerminalOf(symbol.BoolExpr))
	}

	if t.Has(1) && t.At(0).IsNonTerminalOf(symbol.LogicExpr) {
		return p.LogicExpr(ctx, t.At(0))
	}

	panic(rdparser.NewParseError(ctx, "invalid expression"))
}

func (p *Parser) LogicExpr(ctx context.Context, t *rdparser.Tree) bool {
	ctx = rdparser.Trace(ctx, symbol.LogicExpr)

	exprA := p.Expr(ctx, t.At(0).AssertNonTerminalOf(symbol.Expr))
	logicOp := p.LogicOp(ctx, t.At(1).AssertNonTerminalOf(symbol.LogicOp))
	exprB := p.Expr(ctx, t.At(2).AssertNonTerminalOf(symbol.Expr))

	switch logicOp {
	case LogicOpEqu:
		return p.IsEqu(exprA, exprB)
	case LogicOpNotEqu:
		return p.IsNotEqu(exprA, exprB)
	case LogicOpLTEqu:
		return p.IsLTEqu(exprA, exprB)
	case LogicOpGTEqu:
		return p.IsGTEqu(exprA, exprB)
	case LogicOpLT:
		return p.IsLT(exprA, exprB)
	case LogicOpGT:
		return p.IsGT(exprA, exprB)
	}

	panic(rdparser.NewParseError(ctx, "invalid logical op code"))
}

func (p *Parser) LogicOp(ctx context.Context, t *rdparser.Tree) LogicOp {
	ctx = rdparser.Trace(ctx, symbol.LogicOp)

	op := t.At(0).AsTerminal()
	switch op {
	case token.Equ:
		return LogicOpEqu
	case token.NotEquA, token.NotEquB, token.NotEquC:
		return LogicOpNotEqu
	case token.LTEqu:
		return LogicOpLTEqu
	case token.GTEqu:
		return LogicOpGTEqu
	case token.LT:
		return LogicOpLT
	case token.GT:
		return LogicOpGT
	}

	panic(rdparser.NewParseError(ctx, fmt.Sprintf("invalid logical op `%s`", op)))
}

func (p *Parser) Variable(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.Variable)

	varToken := t.At(0).AsTerminal().String()

	varExtract := p.varRegex.FindStringSubmatch(varToken)
	if varExtract == nil {
		panic(rdparser.NewParseError(ctx, fmt.Sprintf("failed to extract variable `%s`", varToken)))
	}
	varName := varExtract[1]

	if rslt, ok := p.Var[varName]; ok {
		return rslt
	}

	panic(rdparser.NewParseError(ctx, fmt.Sprintf("unknown variable `%s`", varName)))
}

func (p *Parser) Number(ctx context.Context, t *rdparser.Tree) float64 {
	ctx = rdparser.Trace(ctx, symbol.Number)

	numToken := t.At(0).AsTerminal().String()
	rslt, err := strconv.ParseFloat(numToken, 64)
	if err != nil {
		panic(rdparser.NewParseError(ctx, fmt.Sprintf("error parsing number `%s`", numToken)))
	}

	return rslt
}

func (p *Parser) IsEqu(x, y float64) bool {
	return math.Abs(x-y) <= p.epsilon
}

func (p *Parser) IsNotEqu(x, y float64) bool {
	return math.Abs(x-y) > p.epsilon
}

func (p *Parser) IsLTEqu(x, y float64) bool {
	return x < y || p.IsEqu(x, y)
}

func (p *Parser) IsGTEqu(x, y float64) bool {
	return x > y || p.IsEqu(x, y)
}

func (p *Parser) IsLT(x, y float64) bool {
	return x < y
}

func (p *Parser) IsGT(x, y float64) bool {
	return x > y
}
