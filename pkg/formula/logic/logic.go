package logic

import "math"

func Equ(x, y, epsilon float64) bool {
	if math.IsInf(x, 1) {
		return math.IsInf(y, 1)
	}

	if math.IsInf(x, -1) {
		return math.IsInf(y, -1)
	}

	if math.IsInf(y, 1) {
		return math.IsInf(x, 1)
	}

	if math.IsInf(y, -1) {
		return math.IsInf(x, -1)
	}

	if math.IsNaN(x) {
		return math.IsNaN(y)
	}

	if math.IsNaN(y) {
		return math.IsNaN(x)
	}

	return math.Abs(x-y) <= epsilon
}

func NotEqu(x, y, epsilon float64) bool {
	return !Equ(x, y, epsilon)
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
