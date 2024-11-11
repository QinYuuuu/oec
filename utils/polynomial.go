package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type Poly struct {
	Coeff []*big.Int // coefficients P(x) = coeff[0] + coeff[1] x + ... + coeff[degree] x^degree ...
}

// NewPoly returns a polynomial P(x) = 0 with capacity degree + 1
func NewPoly(degree int) (Poly, error) {
	if degree < 0 {
		return Poly{}, fmt.Errorf(fmt.Sprintf("degree must be non-negative, got %d", degree))
	}

	coeff := make([]*big.Int, degree+1)

	for i := 0; i < len(coeff); i++ {
		coeff[i] = big.NewInt(0)
	}

	//set the leading coefficient
	//Coeff[len(Coeff) - 1].SetInt64(1)

	return Poly{coeff}, nil
}

// GetDegree returns the degree, ignoring removing leading zeroes
func (poly Poly) GetDegree() int {
	deg := len(poly.Coeff) - 1

	// note: i == 0 is not tested, because even the constant term is zero, we consider it's degree 0
	for i := deg; i > 0; i-- {
		if poly.Coeff[i].Int64() == 0 {
			deg--
		} else {
			break
		}
	}
	return deg
}

// GetLeadingCoefficient returns the coefficient of the highest degree of the variable
func (poly Poly) GetLeadingCoefficient() *big.Int {
	lc := big.NewInt(0)
	lc.Set(poly.Coeff[poly.GetDegree()])

	return lc
}

// GetCoefficient returns Coeff[i]
func (poly Poly) GetCoefficient(i int) (*big.Int, error) {
	if i < 0 || i >= len(poly.Coeff) {
		return big.NewInt(0), errors.New("out of boundary")
	}

	return poly.Coeff[i], nil
}

// SetCoefficient sets the poly.Coeff[i] to ci
func (poly *Poly) SetCoefficient(i int, ci int64) error {
	if i < 0 || i >= len(poly.Coeff) {
		return errors.New("out of boundary")
	}

	poly.Coeff[i].SetInt64(ci)

	return nil
}

// SetCoefficientBig sets the poly.Coeff[i] to ci (a gmp.Int)
func (poly *Poly) SetCoefficientBig(i int, ci *big.Int) error {
	if i < 0 || i >= len(poly.Coeff) {
		return errors.New("out of boundary")
	}

	poly.Coeff[i].Set(ci)

	return nil
}

// Reset sets the coefficients to zeroes
func (poly *Poly) Reset() {
	for i := 0; i < len(poly.Coeff); i++ {
		poly.Coeff[i].SetInt64(0)
	}
}

func (poly *Poly) DeepCopy(other Poly) {
	poly.resetToDegree(other.GetDegree())

	for i := 0; i < other.GetDegree()+1; i++ {
		poly.Coeff[i].Set(other.Coeff[i])
	}
}

// resetToDegree resizes the slice to degree
func (poly *Poly) resetToDegree(degree int) {
	// if we just need to shrink the size
	if degree+1 <= len(poly.Coeff) {
		poly.Coeff = poly.Coeff[:degree+1]
	} else {
		// if we need to grow the slice
		needed := degree + 1 - len(poly.Coeff)
		neededPointers := make([]*big.Int, needed)
		for i := 0; i < len(neededPointers); i++ {
			neededPointers[i] = big.NewInt(0)
		}

		poly.Coeff = append(poly.Coeff, neededPointers...)
	}

	poly.Reset()
}

func (poly Poly) Equal(op Poly) bool {
	if op.GetDegree() != poly.GetDegree() {
		return false
	}

	for i := 0; i <= op.GetDegree(); i++ {
		if op.Coeff[i].Cmp(poly.Coeff[i]) != 0 {
			return false
		}
	}

	return true
}

// IsZero returns if poly == 0
func (poly Poly) IsZero() bool {
	if poly.GetDegree() != 0 {
		return false
	}

	return poly.Coeff[0].Int64() == 0
}

// Rand sets the polynomial coefficients to a pseudo-random number in [0, n)
// WARNING: Rand makes sure that the highest coefficient is not zero
func (poly *Poly) Rand(mod *big.Int) {
	for i := range poly.Coeff {
		poly.Coeff[i], _ = rand.Int(rand.Reader, mod)
	}

	highest := len(poly.Coeff) - 1

	for {
		if poly.Coeff[highest].Int64() == 0 {
			poly.Coeff[highest], _ = rand.Int(rand.Reader, mod)
		} else {
			break
		}
	}

}

func (poly Poly) GetCap() int {
	return len(poly.Coeff)
}

func (poly *Poly) GrowCapTo(cap int) {
	current := poly.GetCap()
	if cap <= current {
		return
	}

	// if we need to grow the slice
	needed := cap - current
	neededPointers := make([]*big.Int, needed)
	for i := 0; i < len(neededPointers); i++ {
		neededPointers[i] = big.NewInt(0)
	}

	poly.Coeff = append(poly.Coeff, neededPointers...)
}

// NewRandPoly returns a randomized polynomial with specified degree
// coefficients are pesudo-random numbers in [0, n)
func NewRandPoly(degree int, n *big.Int) (Poly, error) {
	p, e := NewPoly(degree)
	if e != nil {
		return Poly{}, e
	}

	p.Rand(n)

	return p, nil
}

// NewConstant returns create a constant polynomial P(x) = c
func NewConstant(c int64) Poly {
	zero, err := NewPoly(0)
	if err != nil {
		panic(err.Error())
	}

	zero.Coeff[0] = big.NewInt(c)
	return zero
}

// NewOne creates a constant polynomial P(x) = 1
func NewOne() Poly {
	return NewConstant(1)
}

// NewEmpty creates a constant polynomial P(x) = 0
func NewEmpty() Poly {
	return NewConstant(0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Mod sets poly to poly % p
func (poly *Poly) Mod(p *big.Int) {
	for i := 0; i < len(poly.Coeff); i++ {
		poly.Coeff[i].Mod(poly.Coeff[i], p)
	}
}

// Add sets poly to op1 + op2
func (poly *Poly) Add(op1 Poly, op2 Poly) error {
	// make sure poly is as long as the longest of op1 and op2
	deg1 := op1.GetDegree()
	deg2 := op2.GetDegree()

	if deg1 > deg2 {
		poly.DeepCopy(op1)

	} else {
		poly.DeepCopy(op2)
	}
	for i := 0; i < min(deg1, deg2)+1; i++ {
		poly.Coeff[i].Add(op1.Coeff[i], op2.Coeff[i])
	}
	return nil
}

// AddSelf sets poly to poly + op
func (poly *Poly) AddSelf(op Poly) error {
	var op1 Poly
	op1.DeepCopy(*poly)
	return poly.Add(op1, op)
}

// Sub sets poly to op1 - op2
func (poly *Poly) Sub(op1 Poly, op2 Poly) error {
	// make sure poly is as long as the longest of op1 and op2
	deg1 := op1.GetDegree()
	deg2 := op2.GetDegree()

	if deg1 > deg2 {
		poly.DeepCopy(op1)
	} else {
		poly.DeepCopy(op2)
	}

	for i := 0; i < min(deg1, deg2)+1; i++ {
		poly.Coeff[i].Sub(op1.Coeff[i], op2.Coeff[i])
	}
	poly.Coeff = poly.Coeff[:poly.GetDegree()+1]

	return nil
}

// SubSelf sets poly to poly - op
func (poly *Poly) SubSelf(op Poly) error {
	// make sure poly is as long as the longest of op1 and op2
	deg1 := op.GetDegree()

	poly.GrowCapTo(deg1 + 1)

	for i := 0; i < deg1+1; i++ {
		poly.Coeff[i].Sub(poly.Coeff[i], op.Coeff[i])
	}

	poly.Coeff = poly.Coeff[:poly.GetDegree()+1]

	// FIXME: no need to return error
	return nil
}

// AddMul sets poly to poly + poly2 * k (k being a scalar)
func (poly *Poly) AddMul(poly2 Poly, k *big.Int) {
	for i := 0; i <= poly2.GetDegree(); i++ {
		tmp := new(big.Int).Mul(poly2.Coeff[i], k)
		poly.Coeff[i].Add(poly.Coeff[i], tmp)
	}
}

// Mul set poly to op1 * op2
func (poly *Poly) Mul(op1 Poly, op2 Poly) error {
	deg1 := op1.GetDegree()
	deg2 := op2.GetDegree()

	poly.resetToDegree(deg1 + deg2)

	for i := 0; i <= deg1; i++ {
		for j := 0; j <= deg2; j++ {
			tmp := new(big.Int).Mul(op1.Coeff[i], op2.Coeff[j])
			poly.Coeff[i+j].Add(poly.Coeff[i+j], tmp)
		}
	}

	poly.Coeff = poly.Coeff[:poly.GetDegree()+1]

	return nil
}

// EvalMod returns poly(x) using Horner's rule. If p != nil, returns poly(x) mod p
func (poly Poly) EvalMod(x *big.Int, p *big.Int) *big.Int {
	result := new(big.Int).Set(poly.Coeff[poly.GetDegree()])

	for i := poly.GetDegree(); i >= 1; i-- {
		result.Mul(result, x)
		result.Add(result, poly.Coeff[i-1])
	}

	if p != nil {
		result.Mod(result, p)
	}
	return result
}

// DivMod sets computes q, r such that a = b*q + r.
// This is an implementation of Euclidean division. The complexity is O(n^3)!!
func DivMod(a Poly, b Poly, p *big.Int) (Poly, Poly, error) {
	if b.IsZero() {
		return Poly{}, Poly{}, errors.New("divide by zero")
	}

	var q, r Poly

	q.resetToDegree(0)
	r.DeepCopy(a)

	d := b.GetDegree()
	c := b.GetLeadingCoefficient()

	// cInv = 1/c
	cInv := big.NewInt(0)
	cInv.ModInverse(c, p)

	for r.GetDegree() >= d {
		lc := r.GetLeadingCoefficient()
		s, err := NewPoly(r.GetDegree() - d)
		if err != nil {
			return Poly{}, Poly{}, err
		}

		err = s.SetCoefficientBig(r.GetDegree()-d, lc.Mul(lc, cInv))
		if err != nil {
			return Poly{}, Poly{}, err
		}
		q.AddSelf(s)

		sb := NewEmpty()
		sb.Mul(s, b)

		// deg r reduces by each iteration
		r.SubSelf(sb)

		// modulo p
		q.Mod(p)
		r.Mod(p)
	}

	return q, r, nil
}

func FromVecBig(coeff []*big.Int) Poly {
	if len(coeff) == 0 {
		return NewConstant(0)
	}
	deg := len(coeff) - 1

	poly, err := NewPoly(deg)
	if err != nil {
		panic(err.Error())
	}

	for i := range poly.Coeff {
		poly.Coeff[i] = coeff[i]
	}

	return poly
}

func FromVec(coeff ...int64) Poly {
	if len(coeff) == 0 {
		return NewConstant(0)
	}
	deg := len(coeff) - 1

	poly, err := NewPoly(deg)
	if err != nil {
		panic(err.Error())
	}

	for i := range poly.Coeff {
		poly.Coeff[i].SetInt64(coeff[i])
	}

	return poly
}

func (poly Poly) ToString() string {
	var s = ""

	for i := len(poly.Coeff) - 1; i >= 0; i-- {
		// skip zero coefficients but the constant term
		if i != 0 && poly.Coeff[i].Int64() == 0 {
			continue
		}
		if i > 0 {
			s += fmt.Sprintf("%s x^%d + ", poly.Coeff[i].String(), i)
		} else {
			// constant term
			s += poly.Coeff[i].String()
		}
	}

	return s
}
