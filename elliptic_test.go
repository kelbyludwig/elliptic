package elliptic

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAddition(t *testing.T) {
	curve := NewCurve(big.NewInt(-1), big.NewInt(4), big.NewInt(65535))
	x := NewPoint(big.NewInt(0), big.NewInt(2), curve)
	y := NewPoint(big.NewInt(-1), big.NewInt(-2), curve)
	x.Add(x, y)
	fmt.Printf("Result: (%v, %v)\n", x.X, x.Y)
}
