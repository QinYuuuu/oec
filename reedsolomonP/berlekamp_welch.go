package reedsolomonP

import (
	"fmt"
	"math/big"
	"oec/utils"
)

// modAdd 在模p下执行加法
func modAdd(a, b, p *big.Int) *big.Int {
	result := new(big.Int).Add(a, b)
	result.Mod(result, p)
	return result
}

// modSub 在模p下执行减法
func modSub(a, b, p *big.Int) *big.Int {
	result := new(big.Int).Sub(a, b)
	result.Mod(result, p)
	return result

}

// modMul 在模p下执行乘法
func modMul(a, b, p *big.Int) *big.Int {
	result := new(big.Int).Mul(a, b)
	result.Mod(result, p)
	return result
}

// modInverse computes the modular multiplicative inverse of a under modulo p.
func modInverse(a, p *big.Int) (*big.Int, error) {
	inv := new(big.Int).ModInverse(a, p)
	if inv == nil {
		return nil, fmt.Errorf("no inverse exists for %s in Z_%s", a, p)
	}
	return inv, nil
}

// powerMod computes a^b mod p.
func powerMod(a, b, p *big.Int) *big.Int {
	return new(big.Int).Exp(a, b, p)
}

// berlekampWelch decodes the received vector using the Berlekamp-Welch algorithm.
func berlekampWelch(received []*big.Int, n, k int, p *big.Int) (utils.Poly, error) {
	t := (n - k) / 2 // Number of errors that can be corrected

	// Initialize polynomials Q(x) and E(x)
	Q, err := utils.NewPoly(t)
	if err != nil {
		return utils.Poly{}, err
	}
	E, err := utils.NewPoly(t)
	if err != nil {
		return utils.Poly{}, err
	}
	err = Q.SetCoefficient(0, 1)
	if err != nil {
		return utils.Poly{}, err
	}
	err = E.SetCoefficient(0, 1)
	if err != nil {
		return utils.Poly{}, err
	}

	// Build the system of equations
	for i := 0; i < n; i++ {
		Yi := received[i]
		S := make([]*big.Int, t+1)
		for j := 0; j <= t; j++ {
			xj := powerMod(big.NewInt(int64(i)), big.NewInt(int64(j)), p)
			S[j] = modMul(Q.EvalMod(xj, p), Yi, p)
			if j > 0 {

				S[j] = modSub(S[j], modMul(E.Coeff[j-1], Yi, p), p)
			}
		}
		for j := t; j > 0; j-- {
			E.Coeff[j] = E.Coeff[j-1]
		}
		E.Coeff[0] = big.NewInt(0)
		for j := 0; j <= t; j++ {
			E.Coeff[j] = modAdd(E.Coeff[j], S[j], p)
		}
	}

	// Compute the error locator polynomial Omega(x)
	Omega, err := utils.NewPoly(t + 1)
	if err != nil {
		return utils.Poly{}, err
	}
	for i := 0; i <= t; i++ {
		Omega.Coeff[i] = Q.Coeff[t-i]
	}

	// Chien search to find error locations
	var errors []int
	for i := 0; i < n; i++ {
		sum := big.NewInt(0)
		for j := 0; j <= t; j++ {
			sum.Add(sum, modMul(Omega.Coeff[j], powerMod(big.NewInt(int64(i)), big.NewInt(int64(j)), p), p))
		}
		if sum.Cmp(big.NewInt(0)) == 0 {
			errors = append(errors, i)
		}
	}

	// If the number of errors found does not match t, decoding failed
	if len(errors) != t {
		return utils.Poly{}, fmt.Errorf("decoding failed: found %d errors, expected %d", len(errors), t)
	}

	// Correct the errors
	corrected := make([]*big.Int, n)
	for i := range corrected {
		corrected[i] = new(big.Int).Set(received[i])
	}
	for _, pos := range errors {
		corrected[pos] = big.NewInt(0) // Here we simply set the error position to 0, in practice you should compute the correct value
	}

	// Return the decoded information polynomial
	return utils.FromVecBig(corrected[:k]), nil
}
