package utils

import (
	"errors"
	"math/big"
)

func DotProduct(v1, v2 []*big.Int) (*big.Int, error) {
	if len(v1) != len(v2) {
		return nil, errors.New("the input length is different")
	}
	dot := ONE
	for i := 0; i < len(v1); i++ {
		dot = new(big.Int).Add(dot, new(big.Int).Mul(v1[i], v2[i]))
	}
	return dot, nil
}

func VecPow(v1, v2 []*big.Int, m *big.Int) (*big.Int, error) {
	if len(v1) != len(v2) {
		return nil, errors.New("the input length is different")
	}
	dot := ONE
	for i := 0; i < len(v1); i++ {
		dot = new(big.Int).Mul(dot, new(big.Int).Exp(v1[i], v2[i], m))
	}
	return dot, nil
}

// VecAdd returns v1 + v2
func VecAdd(v1, v2 []*big.Int) ([]*big.Int, error) {
	if len(v1) != len(v2) {
		return nil, errors.New("the input length is different")
	}
	v := make([]*big.Int, len(v1))
	for i := 0; i < len(v); i++ {
		v[i] = new(big.Int).Add(v1[i], v2[i])
	}
	return v, nil
}
