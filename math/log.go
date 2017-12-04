package math

import (
	"github.com/ericlagergren/decimal"
	"github.com/ericlagergren/decimal/internal/arith"
)

// Log sets z to the natural logarithm of x.
func Log(z, x *decimal.Big) *decimal.Big {
	if z.CheckNaNs(x, nil) {
		return z
	}

	if x.Signbit() {
		// ln x is undefined.
		z.Context.Conditions |= decimal.InvalidOperation
		return z.SetNaN(false)
	}

	if x.IsInf(+1) {
		// ln +Inf = +Inf
		return z.SetInf(false)
	}

	if x.Sign() == 0 {
		// ln 0 = -Inf
		return z.SetInf(false)
	}

	if x.IsInt() && !x.IsBig() && x.Uint64() == 1 {
		// ln 1 = 0
		return z.SetMantScale(0, 0)
	}

	t := (x.Scale() - x.Precision()) - 1
	if t < 0 {
		t = -t - 1
	}
	t *= 2

	if arith.Length(arith.Abs(int64(t)))-1 > decimal.MaxScale {
		z.Context.Conditions |= decimal.Overflow | decimal.Inexact | decimal.Rounded
		return z.SetInf(t < 0)
	}

	// Argument reduction:
	// Given
	//    ln(a) = ln(b) + ln(c)
	// Where
	//    a = b * c
	// Given
	//    x = m * 10**n
	// Reduce x (as y) so that
	//    1 <= y <= 10
	// And create p so that
	//    x = y * 10**p
	// Compute
	//    log(y) + p*log(10)

	prec := precision(z)
	y := decimal.WithPrecision(prec).Copy(x).SetScale(x.Precision() - 1)

	var p decimal.Big
	p.SetUint64(arith.Abs(int64(x.Scale() - x.Precision() + 1)))

	// Algorithm is for log(1+x)
	y.Sub(y, one)

	prec += adj
	g := lgen{
		prec: prec,
		pow:  decimal.WithPrecision(prec).Mul(y, y),
		z2:   decimal.WithPrecision(prec).Add(y, two),
		k:    -1,
		t: Term{
			A: decimal.WithPrecision(prec),
			B: decimal.WithPrecision(prec),
		},
	}

	y.Quo(y.Mul(y, two), Lentz(z, &g))
	if p.Sign() != 0 {
		// Avoid the call to log10 if it'll result in 0.
		y.FMA(&p, log10(prec), y)
	}
	return z.Set(y)
}

const adj = 3

type lgen struct {
	prec int
	pow  *decimal.Big // z*z
	z2   *decimal.Big // z+2
	k    int64
	t    Term
}

func (l *lgen) mode() decimal.RoundingMode { return decimal.ToNearestEven }

func (l *lgen) Lentz() (f, Δ, C, D, eps *decimal.Big) {
	f = decimal.WithPrecision(l.prec)
	Δ = decimal.WithPrecision(l.prec)
	C = decimal.WithPrecision(l.prec)
	D = decimal.WithPrecision(l.prec)
	eps = decimal.New(1, l.prec)
	return
}

func (a *lgen) Next() bool { return true }

func (a *lgen) Term() Term {
	// log(z) can be expressed as the following continued fraction:
	//
	//          2z      1^2 * z^2   2^2 * z^2   3^2 * z^2   4^2 * z^2
	//     ----------- ----------- ----------- ----------- -----------
	//      1 * (2+z) - 3 * (2+z) - 5 * (2+z) - 7 * (2+z) - 9 * (2+z) - ···
	//
	// (Cuyt, p 271).
	//
	// References:
	//
	// [Cuyt] - Cuyt, A.; Petersen, V.; Brigette, V.; Waadeland, H.; Jones, W.B.
	// (2008). Handbook of Continued Fractions for Special Functions. Springer
	// Netherlands. https://doi.org/10.1007/978-1-4020-6949-9

	a.k += 2
	if a.k != 1 {
		a.t.A.SetMantScale(-((a.k / 2) * (a.k / 2)), 0)
		a.t.A.Mul(a.t.A, a.pow)
	}
	a.t.B.SetMantScale(a.k, 0)
	a.t.B.Mul(a.t.B, a.z2)
	return a.t
}
