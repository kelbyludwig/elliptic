package elliptic

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAddition(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	x := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	y := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	z := Add(x, y)
	if z.X.Cmp(big.NewInt(20)) != 0 || z.Y.Cmp(big.NewInt(20)) != 0 {
		t.Errorf("(20,20) != (%v, %v)\n", z.X, z.Y)
	}
}

func TestScalarMult(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	x := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	xp := new(Point)
	xp = ScalarMult(x, big.NewInt(9))
	if xp.X.Cmp(big.NewInt(4)) != 0 || xp.Y.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("(4,5) != (%v,%v)\n", xp.X, xp.Y)
	}
}

func TestKeyPair(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	base := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	curve.Base = base
	sec, pub := curve.GenerateKeypair()
	check := ScalarMult(base, sec)
	if !check.Equals(pub) {
		t.Errorf("1. Public key did not match expected scalar multiplcation!\n")
	}

	sec, pub = curve.GenerateKeypair()
	check = ScalarMult(base, sec)
	if !check.Equals(pub) {
		t.Errorf("2. Public key did not match expected scalar multiplcation!\n")
	}

}

func TestOrder(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	base := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	for i := 0; i < 28; i++ {
		p := ScalarMult(base, big.NewInt(int64(i)))
		fmt.Printf("%v: (%v, %v)\n", i, p.X, p.Y)
	}
}

func TestCommutativity(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	p1 := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	p2 := NewPoint(big.NewInt(12), big.NewInt(17), curve)
	r1 := Add(p1, p2)
	r2 := Add(p2, p1)
	fmt.Printf("r1: (%v, %v)\n", r1.X, r1.Y)
	fmt.Printf("r2: (%v, %v)\n", r2.X, r2.Y)
}
