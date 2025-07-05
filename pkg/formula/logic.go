package formula

type LogicOp int

const (
	LogicOpEqu LogicOp = iota
	LogicOpNotEqu
	LogicOpLTEqu
	LogicOpGTEqu
	LogicOpLT
	LogicOpGT
)
