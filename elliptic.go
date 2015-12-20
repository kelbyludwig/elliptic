package elliptic

import "math/big"

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

	if p1.IsIdentity() {
		return new(Point).Set(p2)
	}

	if p2.IsIdentity() {
		return new(Point).Set(p1)
	}

	inverse := new(Point).Set(p2)
	inverse = Inverse(inverse)

	if p1.Equals(inverse) {
		id := Identity(p1.Curve)
		return id
	}

	m := big.NewInt(0)
	if p1.Equals(p2) {
		//Point doubling
		m = new(big.Int).Exp(p1.X, big.NewInt(2), p1.Curve.P)
		m.Mul(big.NewInt(3), m)
		m.Add(m, p1.Curve.A)
		inv := new(big.Int).Mul(big.NewInt(2), p1.Y)
		inv.ModInverse(inv, p1.Curve.P)
		m.Mul(m, inv)
		m = m.Mod(m, p1.Curve.P)
	} else {
		m = new(big.Int).Sub(p1.Y, p2.Y)
		inv := new(big.Int).Sub(p1.X, p2.X)
		inv.ModInverse(inv, p1.Curve.P)
		inv = inv.Mod(inv, p1.Curve.P)
		m = m.Mul(m, inv)
		m = m.Mod(m, p1.Curve.P)
	}

	x := big.NewInt(0)
	y := big.NewInt(0)

	//This is an ugly way of saying: x := m^2 - p1.X - p2.X
	x.Exp(m, big.NewInt(2), p1.Curve.P)
	x.Sub(x, p1.X)
	x.Sub(x, p2.X)
	x.Mod(x, p1.Curve.P)

	y.Sub(p1.X, x)
	y.Mul(m, y)
	y.Sub(y, p1.Y)
	y.Mod(y, p1.Curve.P)

	res := NewPoint(x, y, p1.Curve)
	return res
}

//ScalarMult returns k*pi. If k is zero the identity point is returned
func ScalarMult(pi *Point, k *big.Int) *Point {
	n := new(Point).Set(pi)
	r := Identity(pi.Curve)
	for bit := 0; bit < k.BitLen(); bit++ {
		if k.Bit(bit) == 1 {
			r = Add(r, n)
		}
		n = Add(n, n)
	}
	return r
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
