package elliptic

import (
	"fmt"
	"math/big"
)

var debugAdd bool = false
var debugScalarMult bool = false
var debugSet bool = false
var debugInverse bool = false

//Curve represents an equation of the form: y^2 = x^3 + A*x + B
type Curve struct {
	A    *big.Int
	B    *big.Int
	P    *big.Int
	Base *Point
}

//Point represents a point that (maybe) lies on the curve defined by Curve.
type Point struct {
	X       *big.Int
	Y       *big.Int
	Nonzero bool
	Curve   *Curve
}

//NewCurve initializes a new Curve instance given the supplied parameters.
func NewCurve(a, b, p *big.Int) *Curve {
	curve := &Curve{A: a, B: b, P: p}
	return curve
}

func (c *Curve) equals(ci *Curve) bool {
	if c.A.Cmp(ci.A) == 0 &&
		c.B.Cmp(ci.B) == 0 &&
		c.P.Cmp(ci.P) == 0 {
		return true
	}
	return false
}

//NewPoint initializes a new non-identity point that (maybe) lies on Curve.
func NewPoint(x, y *big.Int, curve *Curve) *Point {
	p := &Point{Nonzero: true, Curve: curve}
	p.X = new(big.Int).Set(x)
	p.Y = new(big.Int).Set(y)
	return p
}

//Identity returns the identity point on Curve c.
func Identity(c *Curve) *Point {
	p := new(Point)
	p.Curve = c
	p.X = new(big.Int)
	p.Y = new(big.Int)
	return p
}

//Add sets p to the result of p1 + p2 and returns p
func Add(p1, p2 *Point) *Point {

	if debugAdd {
		fmt.Printf("[ADD] POINT ADDITION: (%v,%v) + (%v,%v)\n", p1.X, p1.Y, p2.X, p2.Y)
	}

	if p1.IsIdentity() {
		if debugAdd {
			fmt.Printf("[ADD] p1 is identity\n")
		}
		return new(Point).Set(p2)
	}

	if p2.IsIdentity() {
		if debugAdd {
			fmt.Printf("[ADD] p2 is identity\n")
		}
		return new(Point).Set(p1)
	}

	if debugAdd {
		fmt.Printf("[ADD] Setting up p2's inverse\n")
	}
	inverse := new(Point).Set(p2)
	inverse = Inverse(inverse)

	if debugAdd {
		fmt.Printf("[ADD] p2's inverse: (%v, %v)\n", inverse.X, inverse.Y)
	}

	if p1.Equals(inverse) {
		if debugAdd {
			fmt.Printf("[ADD] p1 is equal to the inverse of p2\n")
		}
		id := Identity(p1.Curve)
		return id
	}

	m := big.NewInt(0)
	if p1.Equals(p2) {
		//This is just an ugly way of saying: m := (3 * p1.X ^ 2 + a) / 2 * p1.Y
		m.Exp(p1.X, big.NewInt(2), p1.Curve.P)
		m.Mul(m, big.NewInt(3))
		m.Add(m, p1.Curve.A)
		rhs := new(big.Int).Mul(big.NewInt(2), p1.Y)
		rhs.ModInverse(rhs, p1.Curve.P)
		m.Mul(m, rhs)
	} else {
		if debugAdd {
			fmt.Printf("[ADD] p1 is not equal to p2: (%v, %v) != (%v, %v)\n", p1.X, p1.Y, p2.X, p2.Y)
		}
		//This is just an ugly way of saying: m := (p2.Y - p1.Y) / (p2.X - p1.X)
		//ORIG:m.Sub(p2.Y, p1.Y)
		//ORIG:rhs := new(big.Int).Sub(p2.X, p1.X)
		m.Sub(p1.Y, p2.Y)
		rhs := new(big.Int).Sub(p1.X, p2.X)
		rhs.ModInverse(rhs, p1.Curve.P)
		m.Mul(m, rhs)
	}

	x := big.NewInt(0)
	y := big.NewInt(0)

	//This is an ugly way of saying: x := m^2 - p1.X - p2.X
	x.Exp(m, big.NewInt(2), p1.Curve.P)
	x.Sub(x, p1.X)
	x.Sub(x, p2.X)
	x.Mod(x, p1.Curve.P)

	//This is an ugly way of saying: y := m * (p1.X - x) - p1.Y
	y.Sub(p1.X, x)
	y.Mul(m, y)
	y.Sub(y, p1.Y)
	y.Mod(y, p1.Curve.P)

	res := NewPoint(x, y, p1.Curve)
	return res
}

//ScalarMult returns k*pi. If k is zero the identity point is returned
func ScalarMult(pi *Point, k *big.Int) *Point {
	if big.NewInt(0).Cmp(k) == 0 {
		return Identity(pi.Curve)
	}
	p := new(Point).Set(pi)
	for i := big.NewInt(1); i.Cmp(k) != 0; i.Add(i, big.NewInt(1)) {
		if debugScalarMult {
			fmt.Printf("[MULT] %v of %v: ", new(big.Int).Add(i, big.NewInt(1)), k)
			fmt.Printf("p:(%v, %v) + pi:(%v, %v) = ", p.X, p.Y, pi.X, pi.Y)
		}
		p = Add(p, pi)
		if debugScalarMult {
			fmt.Printf("(%v, %v)\n", p.X, p.Y)
		}
	}
	return p
}

//IsIdentity returns true if p is the identity point.
func (p *Point) IsIdentity() bool {
	if p.Nonzero {
		return false
	}
	return true
}

//Equals returns true if p is the same point as pi
func (p *Point) Equals(pi *Point) bool {
	if p.X.Cmp(pi.X) == 0 && p.Y.Cmp(pi.Y) == 0 {
		return true
	}
	return false
}

//Inverse returns pi's inverse
func Inverse(pi *Point) *Point {
	ny := new(big.Int).Neg(pi.Y)
	p := new(Point).Set(pi)
	p.Y = ny
	return p
}

//Set sets p to the value pi and returns p. This clones the bigints withitn Point.
func (p *Point) Set(pi *Point) *Point {
	p.X = new(big.Int).Set(pi.X)
	p.Y = new(big.Int).Set(pi.Y)
	p.Nonzero = pi.Nonzero
	p.Curve = pi.Curve
	return p
}
