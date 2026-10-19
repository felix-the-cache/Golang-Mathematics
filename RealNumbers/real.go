package main

import (
	"fmt"
	"math"
	"math/big"
)

// Surd is the number A + B*SQRT(N), with A and B stored as exact fractions.
// Example: 2 - SQRT(5) is Surd{A: 2, B: -1, N: 5}.
type Surd struct {
	A, B *big.Rat
	N    int64 // number within the root
}

// NewSurd builds a + bSQRT(n) from whole numbers
func NewSurd(a, b, n int64) Surd {
	return Surd{
		A: big.NewRat(a, 1),
		B: big.NewRat(b, 1),
		N: n,
	}
}

// String prints the surd like "2-1SQRT(5)".
func (s Surd) String() string {
	sign := "+"
	b := new(big.Rat).Set(s.B) // copies B into the new fraction, don't want to change the original when looking to read.
	if b.Sign() < 0 {          // Sign returns either -1, 0, 1
		sign = "-"
		b.Neg(b) // switch to negative sign and copy with the Neg method
	}

	return fmt.Sprintf("%s %s %sSQRT(%d)", s.A.RatString(), sign, b.RatString(), s.N)
}

// float gives a decimal approximation for checking
func (s Surd) Float() float64 {
	a, _ := s.A.Float64()
	b, _ := s.B.Float64()

	return a + b*math.Sqrt(float64(s.N))
}

// Conjugate flips the sign in the middle: a + bn -> a - bn
func (s Surd) Conjugate() Surd {
	return Surd{
		A: s.A,
		B: new(big.Rat).Neg(s.B),
		N: s.N,
	}
}

// Mul Multiplies two surds with the same sqrt(n) using FOIL:
//
//	(a + b√n)(c + d√n) = ac + ad√n + bc√n + bd·n
//	                   = (ac + bd·n) + (ad + bc)√n
func Mul(x, y Surd) Surd {
	ac := new(big.Rat).Mul(x.A, y.A)
	bdn := new(big.Rat).Mul(new(big.Rat).Mul(x.B, y.B), big.NewRat(x.N, 1))
	ad := new(big.Rat).Mul(x.A, y.B)
	bc := new(big.Rat).Mul(x.B, y.A)
	return Surd{
		A: new(big.Rat).Add(ac, bdn),
		B: new(big.Rat).Add(ad, bc),
		N: x.N,
	}
}

func Rationalize(s Surd) (Surd, *big.Rat, error) {
	bottom := Mul(s, s.Conjugate()) // always 0 here
	if bottom.A.Sign() == 0 {
		return Surd{}, nil, fmt.Errorf("Denominator %v is 0", s)
	}
	top := s.Conjugate()
	inv := new(big.Rat).Inv(bottom.A)
	return Surd{
		A: new(big.Rat).Mul(top.A, inv),
		B: new(big.Rat).Mul(top.B, inv),
		N: s.N,
	}, bottom.A, nil
}

func main() {
	// example 1/(2+root5)
	d := NewSurd(2, 1, 5)
	fmt.Println("denominator: ", d)
	fmt.Println("Its conjugate: ", d.Conjugate())

	prod := Mul(d, d.Conjugate())
	fmt.Println("(2+√5)(2-√5)  =   ", prod, " <- the √5 terms cancel, leaving 4 - 5 = -1")

	ans, bottom, _ := Rationalize(d)
	fmt.Println("new denominator:  ", bottom.RatString())
	fmt.Println("1/(2+√5)      =   ", ans)
	fmt.Printf("check: 1/(2+√5) = %.10f,  answer = %.10f\n\n", 1/d.Float(), ans.Float())

	// ---- The other direction: 1 / (2 - √5) ----
	d2 := NewSurd(2, -1, 5) // 2 - √5
	ans2, _, _ := Rationalize(d2)
	fmt.Println("1/(2-√5)      =   ", ans2)
	fmt.Printf("check: 1/(2-√5) = %.10f,  answer = %.10f\n\n", 1/d2.Float(), ans2.Float())

	// ---- More examples ----
	for _, s := range []Surd{NewSurd(3, 1, 2), NewSurd(1, -1, 3), NewSurd(5, 2, 7)} {
		r, bot, _ := Rationalize(s)
		fmt.Printf("1/(%v): bottom becomes %s, answer %v  (check %.6f vs %.6f)\n",
			s, bot.RatString(), r, 1/s.Float(), r.Float())

	}
}
