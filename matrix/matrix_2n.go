package matrix

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Matrix byte[row][col]
type Matrix [][]byte

// newMatrix returns a matrix of zeros.
func newMatrix(rows, cols int) (Matrix, error) {
	if rows <= 0 {
		return nil, errInvalidRowSize
	}
	if cols <= 0 {
		return nil, errInvalidColSize
	}

	m := make([][]byte, rows)
	for i := range m {
		m[i] = make([]byte, cols)
	}
	return m, nil
}

// NewMatrixData initializes a matrix with the given row-major data.
// Note that data is not copied from input.
func newMatrixData(data [][]byte) (Matrix, error) {
	m := Matrix(data)
	err := m.Check()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m Matrix) Check() error {
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
func (m Matrix) String() string {
	rowOut := make([]string, 0, len(m))
	for _, row := range m {
		colOut := make([]string, 0, len(row))
		for _, col := range row {
			colOut = append(colOut, strconv.Itoa(int(col)))
		}
		rowOut = append(rowOut, "["+strings.Join(colOut, ", ")+"]")
	}
	return strings.Join(rowOut, ",\n")
}

// IdentityMatrix returns an identity matrix of the given size.
func identityMatrix(size int) (Matrix, error) {
	m, err := newMatrix(size, size)
	if err != nil {
		return nil, err
	}
	for i := range m {
		m[i][i] = 1
	}
	return m, nil
}

// Multiply multiplies this matrix (the one on the left) by another
// matrix (the one on the right) and returns a new matrix with the result.
// The Mul on elements in the matrix is defined in use and can be over GF(p), GF(2^k).
func (m Matrix) Multiply(right Matrix, Mul func(byte, byte) byte) (Matrix, error) {
	if len(m[0]) != len(right) {
		return nil, fmt.Errorf("columns on left (%d) is different than rows on right (%d)", len(m[0]), len(right))
	}
	result, _ := newMatrix(len(m), len(right[0]))
	for r, row := range result {
		for c := range row {
			var value byte
			for i := range m[0] {
				value ^= Mul(m[r][i], right[i][c])
			}
			result[r][c] = value
		}
	}
	return result, nil
}

// Augment returns the concatenation of this matrix and the matrix on the right.
func (m Matrix) Augment(right Matrix) (Matrix, error) {
	if len(m) != len(right) {
		return nil, errMatrixSize
	}

	result, _ := newMatrix(len(m), len(m[0])+len(right[0]))
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

func (m Matrix) SameSize(n Matrix) error {
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
func (m Matrix) SubMatrix(rmin, cmin, rmax, cmax int) (Matrix, error) {
	result, err := newMatrix(rmax-rmin, cmax-cmin)
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
func (m Matrix) SwapRows(r1, r2 int) error {
	if r1 < 0 || len(m) <= r1 || r2 < 0 || len(m) <= r2 {
		return errInvalidRowSize
	}
	m[r2], m[r1] = m[r1], m[r2]
	return nil
}

// IsSquare will return true if the matrix is square, otherwise false.
func (m Matrix) IsSquare() bool {
	return len(m) == len(m[0])
}

// errSingular is returned if the matrix is singular and cannot be inversed
var errSingular = errors.New("matrix is singular")

// errNotSquare is returned if attempting to inverse a non-square matrix.
var errNotSquare = errors.New("only square matrices can be inverted")

// Invert returns the inverse of this matrix.
// Returns ErrSingular when the matrix is singular and doesn't have an inverse.
// The matrix must be square, otherwise ErrNotSquare is returned.
func (m Matrix) Invert(Mul func(byte, byte) byte, Inverse func(byte) byte) (Matrix, error) {
	if !m.IsSquare() {
		return nil, errNotSquare
	}

	size := len(m)
	work, _ := identityMatrix(size)
	work, _ = m.Augment(work)

	err := work.gaussianElimination(Mul, Inverse)
	if err != nil {
		return nil, err
	}

	return work.SubMatrix(0, size, size, size*2)
}

func (m Matrix) gaussianElimination(Mul func(byte, byte) byte, Inverse func(byte) byte) error {
	rows := len(m)
	columns := len(m[0])
	// Clear out the part below the main diagonal and scale the main
	// diagonal to be 1.
	for r := 0; r < rows; r++ {
		// If the element on the diagonal is 0, find a row below
		// that has a non-zero and swap them.
		if m[r][r] == 0 {
			for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
				if m[rowBelow][r] != 0 {
					err := m.SwapRows(r, rowBelow)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		// If we couldn't find one, the matrix is singular.
		if m[r][r] == 0 {
			return errSingular
		}
		// Scale to 1.
		if m[r][r] != 1 {
			scale := Inverse(m[r][r])
			for c := 0; c < columns; c++ {
				m[r][c] = Mul(m[r][c], scale)
			}
		}
		// Make everything below the 1 be a 0 by subtracting
		// a multiple of it.  (Subtraction and addition are
		// both exclusive or in the Galois field.)
		for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
			if m[rowBelow][r] != 0 {
				scale := m[rowBelow][r]
				for c := 0; c < columns; c++ {
					m[rowBelow][c] ^= Mul(scale, m[r][c])
				}
			}
		}
	}

	// Now clear the part above the main diagonal.
	for d := 0; d < rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m[rowAbove][d] != 0 {
				scale := m[rowAbove][d]
				for c := 0; c < columns; c++ {
					m[rowAbove][c] ^= Mul(scale, m[d][c])
				}

			}
		}
	}
	return nil
}

// Vandermonde creates a Vandermonde matrix, which is guaranteed to have the
// property that any subset of rows that forms a square matrix is invertible.
func Vandermonde(rows, cols int, Exp func(byte, byte) byte) (Matrix, error) {
	result, err := newMatrix(rows, cols)
	if err != nil {
		return nil, err
	}
	for r, row := range result {
		for c := range row {
			result[r][c] = Exp(byte(r+1), byte(c))
		}
	}
	return result, nil
}
