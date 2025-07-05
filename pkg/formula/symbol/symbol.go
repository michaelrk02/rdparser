package symbol

import "github.com/michaelrk02/rdparser"

const (
	Expr    rdparser.NonTerminal = "Expr"
	Exprx   rdparser.NonTerminal = "Expr'"
	Term    rdparser.NonTerminal = "Term"
	Termx   rdparser.NonTerminal = "Term'"
	Factor  rdparser.NonTerminal = "Factor"
	Factorx rdparser.NonTerminal = "Factor'"

	FuncCall rdparser.NonTerminal = "FuncCall"
	FuncName rdparser.NonTerminal = "FuncName"
	FuncArg  rdparser.NonTerminal = "FuncArg"
	FuncArgx rdparser.NonTerminal = "FuncArg'"

	BoolCond   rdparser.NonTerminal = "BoolCond"
	BoolExpr   rdparser.NonTerminal = "BoolExpr"
	BoolExprx  rdparser.NonTerminal = "BoolExpr'"
	BoolTerm   rdparser.NonTerminal = "BoolTerm"
	BoolTermx  rdparser.NonTerminal = "BoolTerm'"
	BoolFactor rdparser.NonTerminal = "BoolFactor"

	LogicExpr rdparser.NonTerminal = "LogicExpr"
	LogicOr   rdparser.NonTerminal = "LogicOr"
	LogicAnd  rdparser.NonTerminal = "LogicAnd"
	LogicNot  rdparser.NonTerminal = "LogicNot"
	LogicOp   rdparser.NonTerminal = "LogicOp"

	Variable rdparser.NonTerminal = "Variable"
	Number   rdparser.NonTerminal = "Number"
)
