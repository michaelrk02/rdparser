# Recursive Descent Parser Library

Written in Golang, based on [rd](https://github.com/shivamMg/rd)

Contains:

- [Mathematical Formula Calculation](#mathematical-formula-calculation)

## Mathematical Formula Calculation

Package: `formula`

### Usage

```
$ go run main.go
  -epsilon float
        use this epsilon value
  -expr string
        expression
  -test string
        load test case from file
```

### Example

```
$ go run main.go -expr "(1.618 + 42) * max((7 > 7 ? 1 : 2), 3)"
130.854
```

### The Context-Free Grammar

```
Expr        -> Term Expr'
Expr'       -> "+" Expr | "-" Expr | NULL
Term        -> Factor Term'
Term'       -> "*" Term | "/" Term | "mod" Term | NULL
Factor      -> "(" Factor' | "-" Factor | Variable | Number | FuncCall
Factor'     ->  BoolCond ")" | Expr ")"

FuncCall    -> FuncName "(" FuncArg ")"
FuncName    -> <function>
FuncArg     -> Expr FuncArg' | NULL
FuncArg'    -> "," FuncArg | NULL

BoolCond    -> BoolExpr "?" Expr ":" Expr
BoolExpr    -> BoolTerm BoolExpr'
BoolExpr'   -> LogicOr BoolExpr | NULL
BoolTerm    -> BoolFactor BoolTerm'
BoolTerm'   -> LogicAnd BoolTerm | NULL
BoolFactor  -> LogicNot BoolFactor | "(" BoolExpr ")" | LogicExpr

LogicExpr   -> Expr LogicOp Expr
LogicOr     -> "||" | "or"
LogicAnd    -> "&&" | "and"
LogicNot    -> "!" | "~" | "not"
LogicOp     -> "==" | "!=" | "~=" | "<>" | "<=" | ">=" | "<" | ">"

Variable    -> <variable>
Number      -> <number>
```
