package qe

import (
	"errors"
)

var ErrNotQuadratic = errors.New("a must not be 0: thats a linear equation")

// Eval returns ax² + bx + c. A root is an x that makes this 0.
func Eval(a, b, c, x float64) float64 {
	return a*x*x + b*x + c
}
