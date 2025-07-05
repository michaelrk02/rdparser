package token

import "github.com/michaelrk02/rdparser"

const (
	Add      rdparser.Terminal = "+"
	Sub      rdparser.Terminal = "-"
	Mul      rdparser.Terminal = "*"
	Div      rdparser.Terminal = "/"
	Mod      rdparser.Terminal = "mod"
	Minus    rdparser.Terminal = "-"
	LParen   rdparser.Terminal = "("
	RParen   rdparser.Terminal = ")"
	LSquare  rdparser.Terminal = "["
	RSquare  rdparser.Terminal = "]"
	Comma    rdparser.Terminal = ","
	Question rdparser.Terminal = "?"
	Colon    rdparser.Terminal = ":"

	Equ     rdparser.Terminal = "=="
	NotEquA rdparser.Terminal = "!="
	NotEquB rdparser.Terminal = "~="
	NotEquC rdparser.Terminal = "<>"
	LTEqu   rdparser.Terminal = "<="
	GTEqu   rdparser.Terminal = ">="
	LT      rdparser.Terminal = "<"
	GT      rdparser.Terminal = ">"

	OrNotation   rdparser.Terminal = "||"
	OrText       rdparser.Terminal = "or"
	AndNotation  rdparser.Terminal = "&&"
	AndText      rdparser.Terminal = "and"
	NotNotationA rdparser.Terminal = "!"
	NotNotationB rdparser.Terminal = "~"
	NotText      rdparser.Terminal = "not"
)

func Dict() []rdparser.Terminal {
	return []rdparser.Terminal{
		Add, Sub, Mul, Div, Mod, Minus, LParen, RParen, LSquare, RSquare, Comma, Question, Colon,
		Equ, NotEquA, NotEquB, NotEquC, LTEqu, GTEqu, LT, GT,
		OrNotation, OrText, AndNotation, AndText, NotNotationA, NotNotationB, NotText,
	}
}
