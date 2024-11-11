package reedsolomonP

import "errors"

// errInvalidRowSize will be returned if attempting to create a matrix with negative or zero row number.
var errInvalidRowSize = errors.New("invalid row size")

// errInvalidColSize will be returned if attempting to create a matrix with negative or zero column number.
var errInvalidColSize = errors.New("invalid column size")

// errColSizeMismatch is returned if the size of matrix columns mismatch.
var errColSizeMismatch = errors.New("column size is not the same for all rows")

// errMatrixSize is returned if matrix dimensions are doesn't match.
var errMatrixSize = errors.New("matrix sizes do not match")

// errNotSquare is returned if matrix dimensions are doesn't match.
var errNotSquare = errors.New("matrix is not square")
