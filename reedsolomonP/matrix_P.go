package reedsolomonP

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type P [][]*big.Int

var BigZero = big.NewInt(0)
var BigOne = big.NewInt(1)

// newMatrix returns a matrix of zeros.
func newMatrixP(rows, cols int) (P, error) {
	if rows <= 0 {
		return nil, errInvalidRowSize
	}
	if cols <= 0 {
		return nil, errInvalidColSize
	}

	m := make([][]*big.Int, rows)
	for i := range m {
		m[i] = make([]*big.Int, cols)
		for j := range m[i] {
			m[i][j] = big.NewInt(0)
		}
	}
	return m, nil
}

// NewMatrixData initializes a matrix with the given row-major data.
// Note that data is not copied from input.
func newMatrixDataP(data [][]*big.Int) (P, error) {
	m := P(data)
	err := m.Check()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m P) Check() error {
	rows := len(m)
	if rows == 0 {
		return errInvalidRowSize
	}
	cols := len(m[0])
	if cols == 0 {
		return errInvalidColSize
	}

	for _, col := range m {
		if len(col) != cols {
			return errColSizeMismatch
		}
	}
	return nil
}

// String returns a human-readable string of the matrix contents.
//
// Example: [[1, 2], [3, 4]]
func (m P) String() string {
	rowOut := make([]string, 0, len(m))
	for _, row := range m {
		colOut := make([]string, 0, len(row))
		for _, col := range row {
			colOut = append(colOut, col.String())
		}
		rowOut = append(rowOut, "["+strings.Join(colOut, ", ")+"]")
	}
	return strings.Join(rowOut, ",\n")
}

// IdentityMatrix returns an identity matrix of the given size.
func identityMatrixP(size int) (P, error) {
	m, err := newMatrixP(size, size)
	if err != nil {
		return nil, err
	}
	for i := range m {
		m[i][i] = new(big.Int).SetInt64(1)
	}
	return m, nil
}

// Multiply multiplies this matrix (the one on the left) by another
// matrix (the one on the right) and returns a new matrix with the result.
func (m P) Multiply(right P, p *big.Int) (P, error) {
	if len(m[0]) != len(right) {
		return nil, fmt.Errorf("columns on left (%d) is different than rows on right (%d)", len(m[0]), len(right))
	}
	result, _ := newMatrixP(len(m), len(right[0]))
	for r, row := range result {
		for c := range row {
			value := new(big.Int).SetInt64(0)
			for i := range m[0] {
				tmp := new(big.Int).Mul(m[r][i], right[i][c])
				value = new(big.Int).Add(value, tmp)
			}
			result[r][c] = new(big.Int).Mod(value, p)
		}
	}
	return result, nil
}

// Augment returns the concatenation of this matrix and the matrix on the right.
func (m P) Augment(right P) (P, error) {
	if len(m) != len(right) {
		return nil, errMatrixSize
	}

	result, _ := newMatrixP(len(m), len(m[0])+len(right[0]))
	for r, row := range m {
		for c := range row {
			result[r][c] = m[r][c]
		}
		cols := len(m[0])
		for c := range right[0] {
			result[r][cols+c] = right[r][c]
		}
	}
	return result, nil
}

func (m P) SameSize(n P) error {
	if len(m) != len(n) {
		return errMatrixSize
	}
	for i := range m {
		if len(m[i]) != len(n[i]) {
			return errMatrixSize
		}
	}
	return nil
}

// SubMatrix returns a part of this matrix. Data is copied.
func (m P) SubMatrix(rmin, cmin, rmax, cmax int) (P, error) {
	result, err := newMatrixP(rmax-rmin, cmax-cmin)
	if err != nil {
		return nil, err
	}
	// OPTME: If used heavily, use copy function to copy slice
	for r := rmin; r < rmax; r++ {
		for c := cmin; c < cmax; c++ {
			result[r-rmin][c-cmin] = m[r][c]
		}
	}
	return result, nil
}

// SwapRows Exchanges two rows in the matrix.
func (m P) SwapRows(r1, r2 int) error {
	if r1 < 0 || len(m) <= r1 || r2 < 0 || len(m) <= r2 {
		return errInvalidRowSize
	}
	m[r2], m[r1] = m[r1], m[r2]
	return nil
}

// IsSquare will return true if the matrix is square, otherwise false.
func (m P) IsSquare() bool {
	return len(m) == len(m[0])
}

// Invert returns the inverse of this matrix.
// Returns ErrSingular when the matrix is singular and doesn't have an inverse.
// The matrix must be square, otherwise ErrNotSquare is returned.
func (m P) Invert(p *big.Int) (P, error) {
	if !m.IsSquare() {
		return nil, errNotSquare
	}

	size := len(m)
	work, _ := identityMatrixP(size)
	work, _ = m.Augment(work)

	err := work.gaussianElimination(p)
	if err != nil {
		return nil, err
	}

	return work.SubMatrix(0, size, size, size*2)
}

func (m P) gaussianElimination(p *big.Int) error {
	rows := len(m)
	columns := len(m[0])
	// Clear out the part below the main diagonal and scale the main
	// diagonal to be 1.
	for r := 0; r < rows; r++ {
		// If the element on the diagonal is 0, find a row below
		// that has a non-zero and swap them.
		if m[r][r].Cmp(BigZero) == 0 {
			for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
				if m[rowBelow][r].Cmp(BigZero) != 0 {
					err := m.SwapRows(r, rowBelow)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		// If we couldn't find one, the matrix is singular.
		if m[r][r].Cmp(BigZero) == 0 {
			return errors.New("matrix is singular")
		}
		// Scale to 1.
		if m[r][r].Cmp(BigOne) != 0 {
			scale := new(big.Int).ModInverse(m[r][r], p)
			for c := 0; c < columns; c++ {
				m[r][c] = new(big.Int).Mod(new(big.Int).Mul(m[r][c], scale), p)
			}
		}
		// Make everything below the 1 be a 0 by subtracting
		// a multiple of it.  (Subtraction and addition are
		// both exclusive or in the Galois field.)
		for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
			if m[rowBelow][r].Cmp(BigOne) != 0 {
				scale := m[rowBelow][r]
				for c := 0; c < columns; c++ {
					tmp := new(big.Int).Mul(scale, m[r][c])
					m[rowBelow][c] = new(big.Int).Mod(new(big.Int).Add(m[rowBelow][c], tmp), p)
				}
			}
		}
	}

	// Now clear the part above the main diagonal.
	for d := 0; d < rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m[rowAbove][d].Cmp(BigZero) != 0 {
				scale := m[rowAbove][d]
				for c := 0; c < columns; c++ {
					tmp := new(big.Int).Mul(scale, m[d][c])
					m[rowAbove][c] = new(big.Int).Mod(new(big.Int).Add(m[rowAbove][c], tmp), p)
				}

			}
		}
	}
	return nil
}

// VandermondeP creates a Vandermonde matrix, which is guaranteed to have the
// property that any subset of rows that forms a square matrix is invertible.
func VandermondeP(rows, cols int, p *big.Int) (P, error) {
	result, err := newMatrixP(rows, cols)
	if err != nil {
		return nil, err
	}
	for r, row := range result {
		for c := range row {
			result[r][c] = new(big.Int).Exp(new(big.Int).SetInt64(int64(r+1)), new(big.Int).SetInt64(int64(c)), p)
		}
	}
	return result, nil
}
