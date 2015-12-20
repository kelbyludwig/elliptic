package elliptic

import (
	"fmt"
	"math/big"
	"testing"
)

func TestSet(t *testing.T) {
	curve := NewCurve(big.NewInt(9), big.NewInt(17), big.NewInt(23))
	x := NewPoint(big.NewInt(1), big.NewInt(18), curve)
	y := NewPoint(big.NewInt(11), big.NewInt(13), curve)
	ybefore := fmt.Sprintf("(%v,%v)", y.X, y.Y)
	x.Set(y)
	xafter := fmt.Sprintf("(%v,%v)", x.X, x.Y)
	yafter := fmt.Sprintf("(%v,%v)", y.X, y.Y)
	if xafter != ybefore || ybefore != yafter {
		t.Errorf("TestSet failed!\n")
	}
}

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
	x = ScalarMult(x, big.NewInt(2))
	fmt.Printf("2p: (%v, %v)\n", x.X, x.Y)
	t.Errorf("DONE\n")
}
