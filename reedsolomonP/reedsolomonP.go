package reedsolomonP

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
)

// RSGFp online-error correction algorithm in modulo P Field
// k required pieces and n total pieces.
type RSGFp struct {
	k         int
	n         int
	encMatrix P
	p         *big.Int
}

func NewRSGFp(k, n int, p *big.Int) (*RSGFp, error) {
	if k <= 0 || n <= 0 || k > n {
		return nil, errors.New("requires 1 <= k <= n <= 256")
	}

	encMatrix, err := VandermondeP(n, k, p)
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

// Encode will take input data and encode to the total number of pieces n this
// *FEC is configured for.
//
// The input data must be a multiple of the required number of pieces k.
// Padding to this multiple is up to the caller.
func (oec *RSGFp) Encode(input []*big.Int) ([]Share, error) {
	if len(input) < oec.k {
		return nil, ErrTooFewShards
	}
	output := make([]Share, oec.n)

	for i := 0; i < oec.n; i++ {
		fecBuf := new(big.Int).SetInt64(0)
		for j := 0; j < oec.k; j++ {
			fecBuf = new(big.Int).Add(fecBuf, new(big.Int).Mul(input[j], oec.encMatrix[i][j]))
			fecBuf.Mod(fecBuf, oec.p)
		}

		output[i] = Share{
			Number: i,
			Data:   fecBuf,
		}
	}
	return output, nil
}

// Rebuild will take a list of corrected shares (pieces) and a callback output.
// output will be called k times ((*FEC).Required() times) with 1/k of the
// original data each time and the index of that data piece.
// Decode is usually preferred.
//
// Note that the data is not necessarily sent to output ordered by the piece
// number.
//
// Note that the byte slices in Shares passed to output may be reused when
// output returns.
//
// Rebuild assumes that you have already called Correct or did not need to.
func (oec *RSGFp) Rebuild(shares []Share, output func(Share)) error {
	k := oec.k
	n := oec.n
	encMatrix := oec.encMatrix

	if len(shares) < k {
		return ErrTooFewShards
	}

	sort.Sort(byNumber(shares))
	fmt.Println(shares)
	mDec := make([][]*big.Int, k)
	for i := 0; i < k; i++ {
		mDec[i] = make([]*big.Int, k)
	}
	indexes := make([]int, k)
	sharesv := make([]*big.Int, k)

	sharesBIter := 0
	sharesEIter := len(shares) - 1

	for i := 0; i < k; i++ {
		var shareID int
		var shareData *big.Int
		share := shares[sharesBIter]
		if share.Number == i {
			shareID = share.Number
			shareData = share.Data
			sharesBIter++
		} else {
			share1 := shares[sharesEIter]
			shareID = share1.Number
			shareData = share1.Data
			sharesEIter--
		}

		if shareID >= n {
			return fmt.Errorf("invalid share id: %d", shareID)
		}

		if shareID < k {
			mDec[i][i] = BigOne
			if output != nil {
				output(Share{
					Number: shareID,
					Data:   shareData})
			}
		} else {
			copy(mDec[i], encMatrix[shareID][:k])
		}

		sharesv[i] = shareData
		indexes[i] = shareID
	}

	// Solve the system of linear equations to find the original data
	var buf *big.Int
	for i := 0; i < k; i++ {
		if indexes[i] >= k {
			// Calculate the data for the missing share
			buf = big.NewInt(0)
			for col := 0; col < k; col++ {
				product := new(big.Int).Mul(sharesv[col], mDec[i][col])
				buf.Add(buf, product)
				buf.Mod(buf, oec.p)
			}
			if output != nil {
				output(Share{
					Number: i,
					Data:   buf})
			}
		}
	}

	return nil
}
