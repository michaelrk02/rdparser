package formula

import (
	"context"
	"fmt"
	"math"

	"github.com/michaelrk02/rdparser"
)

type Function func(ctx context.Context, args []float64) float64

type Library interface {
	Resolve(funcName string) (Function, bool)
}

type StdLibrary struct {
	ref map[string]Function
}

func NewStdLibrary() *StdLibrary {
	lib := &StdLibrary{
		ref: make(map[string]Function),
	}

	lib.ref["pow"] = lib.Pow
	lib.ref["round"] = lib.Round
	lib.ref["min"] = lib.Min
	lib.ref["max"] = lib.Max
	lib.ref["sum"] = lib.Sum
	lib.ref["avg"] = lib.Avg
	lib.ref["average"] = lib.Avg

	return lib
}

func (lib *StdLibrary) Resolve(funcName string) (Function, bool) {
	if fn, ok := lib.ref[funcName]; ok {
		return fn, true
	}
	return nil, false
}

func (lib *StdLibrary) Pow(ctx context.Context, args []float64) float64 {
	Validate(ctx, "pow", args).ArgLength(2)

	return math.Pow(args[0], args[1])
}

func (lib *StdLibrary) Round(ctx context.Context, args []float64) float64 {
	Validate(ctx, "round", args).ArgLength(2)

	fac := math.Pow10(int(args[1]))
	return math.Round(args[0]*fac) / fac
}

func (lib *StdLibrary) Min(ctx context.Context, args []float64) float64 {
	x := math.Inf(1)
	for _, n := range args {
		x = math.Min(x, n)
	}
	return x
}

func (lib *StdLibrary) Max(ctx context.Context, args []float64) float64 {
	x := math.Inf(-1)
	for _, n := range args {
		x = math.Max(x, n)
	}
	return x
}

func (lib *StdLibrary) Sum(ctx context.Context, args []float64) float64 {
	x := 0.0
	for _, n := range args {
		x = x + n
	}
	return x
}

func (lib *StdLibrary) Avg(ctx context.Context, args []float64) float64 {
	sum := 0.0
	for _, n := range args {
		sum = sum + n
	}
	return sum / float64(len(args))
}

type Validator struct {
	Ctx      context.Context
	FuncName string
	Args     []float64
}

func Validate(ctx context.Context, funcName string, args []float64) *Validator {
	return &Validator{
		Ctx:      ctx,
		FuncName: funcName,
		Args:     args,
	}
}

func (v *Validator) Error(msg string) error {
	return rdparser.NewParseError(v.Ctx, fmt.Sprintf("[%s] - %s", v.FuncName, msg))
}

func (v *Validator) Rule(validatorFunc func(v *Validator)) {
	validatorFunc(v)
}

func (v *Validator) ArgLength(n int) {
	if len(v.Args) != n {
		panic(v.Error(fmt.Sprintf("expected %d arguments, got %d instead", n, len(v.Args))))
	}
}

func (v *Validator) ArgMinLength(n int) {
	if len(v.Args) < n {
		panic(v.Error(fmt.Sprintf("expected at least %d arguments, got %d instead", n, len(v.Args))))
	}
}
