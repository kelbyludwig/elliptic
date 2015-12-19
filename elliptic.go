package elliptic

import (
	"math/big"
)

//TODO: Define Curve equality for point addition error checking!

//Curve represents an equation of the form: y^2 = x^3 + A*x + B
type Curve struct {
	A    *big.Int
	B    *big.Int
	P    *big.Int
	Base *Point
}

//Point represents a point that (maybe) lies on the curve defined by Curve.
type Point struct {
	X *big.Int
	Y *big.Int
	//nonzero allows for a zero value Point struct to act as the identitiy point
	Nonzero bool
	Curve   *Curve
}

//NewCurve initializes a new Curve instance given the supplied parameters.
func NewCurve(a, b, p *big.Int) *Curve {
	curve := &Curve{A: a, B: b, P: p}
	return curve
}

//NewPoint initializes a new non-identity point that (maybe) lies on Curve.
func NewPoint(x, y *big.Int, curve *Curve) *Point {
	p := &Point{X: x, Y: y, Nonzero: true, Curve: curve}
	return p
}

//Identity returns the identity point on Curve c.
func Identity(c *Curve) *Point {
	var p *Point
	p.Curve = c
	return p
}

//Add sets p1 to the result of p1 + p2 and returns p1
func (p1 *Point) Add(p2 *Point) (p3 *Point) {
	if p1.IsIdentity() {
		return p1.Set(p2)
	}
	if p2.IsIdentity() {
		return p1
	}

	var inverse *Point
	inverse.Set(p2)
	inverse.Inverse(inverse)

	if p1.Equals(inverse) {
		return Identity(p1.Curve)
	}

	var m *big.Int
	if p1.Equals(p2) {
		//This is just an ugly way of saying: m := (3 * p1.X ^ 2 + a) / 2 * p1.Y
		m.Exp(p1.X, big.NewInt(2), p1.Curve.P)
		m.Mul(m, big.NewInt(3))
		m.Add(m, p1.Curve.A)
		var rhs *big.Int
		rhs.Mul(big.NewInt(2), p1.Y)
		m.Div(m, rhs)
	} else {
		//This is just an ugly way of saying: m := (p2.Y - p1.Y) / (p2.X - p1.X)
		m.Sub(p2.Y, p1.Y)
		var rhs *big.Int
		rhs.Sub(p2.X, p1.X)
		m.Div(m, rhs)
	}

	var x *big.Int
	var y *big.Int

	//This is an ugly way of saying: x := m^2 - p1.X - p2.X
	x.Exp(m, big.NewInt(2), p1.Curve.P)
	x.Sub(x, p1.X)
	x.Sub(x, p2.X)

	//This is an ugly way of saying: y := m * (p1.X - x) - p1.Y
	y.Sub(p1.X, x)
	y.Mul(m, y)
	y.Sub(y, p1.Y)

	return NewPoint(x, y, p1.Curve)
}

//IsIdentity returns true if p is the identity point.
func (p *Point) IsIdentity() bool {
	if p.Nonzero {
		return false
	}
	return true
}

func (p *Point) Equals(pi *Point) bool {
	if p.X.Cmp(pi.X) == 0 && p.Y.Cmp(pi.Y) == 0 {
		return true
	}
	return false
}

//Inverse sets p to pi's inverse and returns p
func (p *Point) Inverse(pi *Point) *Point {
	p.Set(pi)
	p.Y.Neg(pi.Y)
	p.Y.Mod(pi.Y, pi.Curve.P)
	return p
}

//Set sets p to the pi and returns p
func (p *Point) Set(pi *Point) *Point {
	p.X = pi.X
	p.Y = pi.Y
	p.Nonzero = pi.Nonzero
	return p
}
