package reedsolomonP

import (
	"math/big"
	"testing"
)

func TestNewOEC(t *testing.T) {
	n := 4
	f := 1
	p := new(big.Int).SetInt64(13)
	_, err := NewOECGFp(f+1, n, p)
	if err != nil {
		t.Fatal(err)
	}
}
