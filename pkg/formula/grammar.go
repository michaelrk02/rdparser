package formula

import (
	"context"
	"fmt"
	"regexp"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula/pattern"
	"github.com/michaelrk02/rdparser/pkg/formula/symbol"
	"github.com/michaelrk02/rdparser/pkg/formula/token"
)

/*
	Grammar:

	Expr		-> Term Expr'
	Expr'		-> "+" Expr | "-" Expr | NULL
	Term		-> Factor Term'
	Term'		-> "*" Term | "/" Term | "mod" Term | NULL
	Factor		-> "(" Factor' | "-" Factor | Variable | Number | FuncCall
	Factor'		->  BoolCond ")" | Expr ")"

	FuncCall	-> FuncName "(" FuncArg ")"
	FuncName	-> <function>
	FuncArg		-> Expr FuncArg' | NULL
	FuncArg'	-> "," FuncArg | NULL

	BoolCond	-> BoolExpr "?" Expr ":" Expr
	BoolExpr	-> BoolTerm BoolExpr'
	BoolExpr'	-> LogicOr BoolExpr | NULL
	BoolTerm	-> BoolFactor BoolTerm'
	BoolTerm'	-> LogicAnd BoolTerm | NULL
	BoolFactor	-> LogicNot BoolFactor | "(" BoolExpr ")" | LogicExpr

	LogicExpr	-> Expr LogicOp Expr
	LogicOr		-> "||" | "or"
	LogicAnd	-> "&&" | "and"
	LogicNot	-> "!" | "~" | "not"
	LogicOp		-> "==" | "!=" | "~=" | "<>" | "<=" | ">=" | "<" | ">"

	Variable	-> <variable>
	Number		-> <number>
*/

type Grammar struct {
	FunctionPattern *regexp.Regexp
	VariablePattern *regexp.Regexp
	NumberPattern   *regexp.Regexp
}

func NewGrammar() *Grammar {
	return &Grammar{
		FunctionPattern: regexp.MustCompile(fmt.Sprintf(`^%s$`, pattern.Function)),
		VariablePattern: regexp.MustCompile(fmt.Sprintf(`^%s$`, pattern.Variable)),
		NumberPattern:   regexp.MustCompile(fmt.Sprintf(`^%s$`, pattern.Number)),
	}
}

func (g *Grammar) BuildParseTree(ctx context.Context, b *rdparser.Builder) (err error) {
	defer rdparser.Catch(rdparser.ErrCompile, &err)

	ok := g.Expr(ctx, b)
	if !ok {
		err = rdparser.NewSyntaxError(ctx, b)
		return
	}

	return
}

func (g *Grammar) Expr(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Expr).Exit(&ok)

	return g.Term(ctx, b) && g.Exprx(ctx, b)
}

func (g *Grammar) Exprx(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Exprx).Exit(&ok)

	if b.Match(token.Add) {
		return g.Expr(ctx, b)
	}

	if b.Match(token.Sub) {
		return g.Expr(ctx, b)
	}

	return true
}

func (g *Grammar) Term(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Term).Exit(&ok)

	return g.Factor(ctx, b) && g.Termx(ctx, b)
}

func (g *Grammar) Termx(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Termx).Exit(&ok)

	if b.Match(token.Mul) {
		return g.Term(ctx, b)
	}

	if b.Match(token.Div) {
		return g.Term(ctx, b)
	}

	if b.Match(token.Mod) {
		return g.Term(ctx, b)
	}

	return true
}

func (g *Grammar) Factor(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Factor).Exit(&ok)

	if b.Match(token.LParen) {
		return g.Factorx(ctx, b)
	}

	if b.Match(token.Minus) {
		return g.Factor(ctx, b)
	}

	if g.Variable(ctx, b) || g.Number(ctx, b) {
		return true
	}

	return g.FuncCall(ctx, b)
}

func (g *Grammar) Factorx(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Factorx).Exit(&ok)

	if g.BoolCond(ctx, b) {
		return b.Match(token.RParen)
	}

	if g.Expr(ctx, b) {
		return b.Match(token.RParen)
	}

	return false
}

func (g *Grammar) FuncCall(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.FuncCall).Exit(&ok)

	return g.FuncName(ctx, b) && b.Match(token.LParen) && g.FuncArg(ctx, b) && b.Match(token.RParen)
}

func (g *Grammar) FuncName(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.FuncName).Exit(&ok)

	tok, ok := b.Next()
	if !ok {
		return false
	}

	if g.FunctionPattern.MatchString(tok.(rdparser.Terminal).String()) {
		b.Add(tok)
		return true
	}

	return false
}

func (g *Grammar) FuncArg(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.FuncArg).Exit(&ok)

	return g.Expr(ctx, b) && g.FuncArgx(ctx, b)
}

func (g *Grammar) FuncArgx(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.FuncArgx).Exit(&ok)

	if b.Match(token.Comma) {
		return g.FuncArg(ctx, b)
	}

	return true
}

func (g *Grammar) BoolCond(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.BoolCond).Exit(&ok)

	return g.BoolExpr(ctx, b) && b.Match(token.Question) && g.Expr(ctx, b) && b.Match(token.Colon) && g.Expr(ctx, b)
}

func (g *Grammar) BoolExpr(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.BoolExpr).Exit(&ok)

	return g.BoolTerm(ctx, b) && g.BoolExprx(ctx, b)
}

func (g *Grammar) BoolExprx(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.BoolExprx).Exit(&ok)

	if g.LogicOr(ctx, b) {
		return g.BoolExpr(ctx, b)
	}

	return true
}

func (g *Grammar) BoolTerm(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.BoolTerm).Exit(&ok)

	return g.BoolFactor(ctx, b) && g.BoolTermx(ctx, b)
}

func (g *Grammar) BoolTermx(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.BoolTermx).Exit(&ok)

	if g.LogicAnd(ctx, b) {
		return g.BoolTerm(ctx, b)
	}

	return true
}

func (g *Grammar) BoolFactor(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.BoolFactor).Exit(&ok)

	if g.LogicNot(ctx, b) {
		return g.BoolFactor(ctx, b)
	}

	if b.Match(token.LParen) && g.BoolExpr(ctx, b) && b.Match(token.RParen) {
		return true
	}
	b.Backtrack()

	return g.LogicExpr(ctx, b)
}

func (g *Grammar) LogicExpr(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.LogicExpr).Exit(&ok)

	return g.Expr(ctx, b) && g.LogicOp(ctx, b) && g.Expr(ctx, b)
}

func (g *Grammar) LogicOr(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.LogicOr).Exit(&ok)

	return b.Match(token.OrNotation) || b.Match(token.OrText)
}

func (g *Grammar) LogicAnd(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.LogicAnd).Exit(&ok)

	return b.Match(token.AndNotation) || b.Match(token.AndText)
}

func (g *Grammar) LogicNot(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.LogicNot).Exit(&ok)

	return b.Match(token.NotNotationA) || b.Match(token.NotNotationB) || b.Match(token.NotText)
}

func (g *Grammar) LogicOp(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.LogicOp).Exit(&ok)

	tok, ok := b.Next()
	if !ok {
		return false
	}

	if g.IsLogicOp(tok.(rdparser.Terminal)) {
		b.Add(tok)
		return true
	}

	return false
}

func (g *Grammar) Variable(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Variable).Exit(&ok)

	tok, ok := b.Next()
	if !ok {
		return false
	}

	if g.VariablePattern.MatchString(tok.(rdparser.Terminal).String()) {
		b.Add(tok)
		return true
	}

	return false
}

func (g *Grammar) Number(ctx context.Context, b *rdparser.Builder) (ok bool) {
	defer b.Enter(&ctx, symbol.Number).Exit(&ok)

	tok, ok := b.Next()
	if !ok {
		return false
	}

	if g.NumberPattern.MatchString(tok.(rdparser.Terminal).String()) {
		b.Add(tok)
		return true
	}

	return false
}

func (g *Grammar) IsLogicOp(tok rdparser.Terminal) bool {
	return tok == token.Equ ||
		tok == token.NotEquA ||
		tok == token.NotEquB ||
		tok == token.NotEquC ||
		tok == token.LTEqu ||
		tok == token.GTEqu ||
		tok == token.LT ||
		tok == token.GT
}
