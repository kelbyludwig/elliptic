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
	t.Errorf("TestAddition: (%v,%v)\n", z.X, z.Y)
}

func TestScalarMult(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	x := NewPoint(big.NewInt(16), big.NewInt(5), curve)
	fmt.Printf("Original p : (%v, %v)\n", x.X, x.Y)
	for i := int64(0); i < 10; i++ {
		xp := ScalarMult(x, big.NewInt(i))
		fmt.Printf("%vp: (%v, %v)\n", i, xp.X, xp.Y)
	}
}
