package field

import (
	"math/big"
)

// GFp represents a finite field modulo a prime number
type GFp struct {
	p *big.Int // 模数 P
}

// NewGFp returns a new finite field modulo a prime number
func NewGFp(p *big.Int) *GFp {
	return &GFp{
		p: p,
	}
}

// Add performs addition modulo P
func (f *GFp) Add(a, b []byte) []byte {
	x := new(big.Int).SetBytes(a)
	y := new(big.Int).SetBytes(b)
	result := new(big.Int).Add(x, y)
	result.Mod(result, f.p)
	return result.Bytes()
}

// Subtract performs subtraction modulo P
func (f *GFp) Subtract(a, b []byte) []byte {
	x := new(big.Int).SetBytes(a)
	y := new(big.Int).SetBytes(b)
	result := new(big.Int).Sub(x, y)
	result.Mod(result, f.p)
	return result.Bytes()
}

// Multiply 在模 P 下执行乘法
func (f *GFp) Multiply(a, b []byte) []byte {
	x := new(big.Int).SetBytes(a)
	y := new(big.Int).SetBytes(b)
	result := new(big.Int).Mul(x, y)
	result.Mod(result, f.p)
	return result.Bytes()
}

// Divide 在模 P 下执行除法（乘以逆元）
func (f *GFp) Divide(a, b []byte) []byte {
	x := new(big.Int).SetBytes(a)
	y := new(big.Int).SetBytes(b)
	inverse := new(big.Int).ModInverse(y, f.p)
	result := new(big.Int).Mul(x, inverse)
	result.Mod(result, f.p)
	return result.Bytes()
}

// Inverse 计算元素在模 P 下的逆元
func (f *GFp) Inverse(a []byte) []byte {
	x := new(big.Int).SetBytes(a)
	return new(big.Int).ModInverse(x, f.p).Bytes()
}

// Zero 返回模 P 下的零元素
func (f *GFp) Zero() []byte {
	return big.NewInt(0).Bytes()
}

// One 返回模 P 下的单位元素
func (f *GFp) One() []byte {
	return big.NewInt(1).Bytes()
}
