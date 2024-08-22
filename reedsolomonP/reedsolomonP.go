package reedsolomonP

import (
	"errors"
	"fmt"
	"math/big"
	"oec/matrix"
)

// RSGFp online-error correction algorithm in modulo P Field
// k required pieces and n total pieces.
type RSGFp struct {
	k         int
	n         int
	encMatrix matrix.P
	p         *big.Int
}

func NewRSGFp(k, n int, p *big.Int) (*RSGFp, error) {
	if k <= 0 || n <= 0 || k > n {
		return nil, errors.New("requires 1 <= k <= n <= 256")
	}

	encMatrix, err := matrix.VandermondeP(n, k, p)
	if err != nil {
		return nil, err
	}
	fmt.Println(encMatrix.String())
	return &RSGFp{
		k:         k,
		n:         n,
		encMatrix: encMatrix,
		p:         p,
	}, nil
}

// ErrTooFewShards is returned if too few shards where given to
// Encode/Verify/Reconstruct/Update. It will also be returned from Reconstruct
// if there were too few shards to reconstruct the missing data.
var ErrTooFewShards = errors.New("too few shards given")

func (oec *RSGFp) Encode(input [][]byte) ([][]byte, error) {
	if len(input) != oec.k {
		return nil, ErrTooFewShards
	}
	output := make([][]byte, oec.n)
	for i := 0; i < oec.n; i++ {
		output[i] = make([]byte, oec.k)
	}

	return output, nil
}
