package logic

import "math"

func Equ(x, y, epsilon float64) bool {
	return math.Abs(x-y) <= epsilon
}

func NotEqu(x, y, epsilon float64) bool {
	return math.Abs(x-y) > epsilon
}

func LTEqu(x, y, epsilon float64) bool {
	return x < y || Equ(x, y, epsilon)
}

func GTEqu(x, y, epsilon float64) bool {
	return x > y || Equ(x, y, epsilon)
}

func LT(x, y float64) bool {
	return x < y
}

func GT(x, y float64) bool {
	return x > y
}
