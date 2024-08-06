package matrix

import (
	"math/big"
	"testing"
)

// TestNewMatrix - Tests validate the result for invalid input and the allocations made by newMatrix method.
func TestNewMatrixP(t *testing.T) {
	testCases := []struct {
		rows    int
		columns int

		// flag to indicate whether the test should pass.
		shouldPass     bool
		expectedResult P
		expectedErr    error
	}{
		// Test case - 1.
		// Test case with a negative row size.
		{-1, 10, false, nil, errInvalidRowSize},
		// Test case - 2.
		// Test case with a negative column size.
		{10, -1, false, nil, errInvalidColSize},
		// Test case - 3.
		// Test case with negative value for both row and column size.
		{-1, -1, false, nil, errInvalidRowSize},
		// Test case - 4.
		// Test case with 0 value for row size.
		{0, 10, false, nil, errInvalidRowSize},
		// Test case - 5.
		// Test case with 0 value for column size.
		{-1, 0, false, nil, errInvalidRowSize},
		// Test case - 6.
		// Test case with 0 value for both row and column size.
		{0, 0, false, nil, errInvalidRowSize},
	}
	for i, testCase := range testCases {
		actualResult, actualErr := newMatrix(testCase.rows, testCase.columns)
		if actualErr != nil && testCase.shouldPass {
			t.Errorf("Test %d: Expected to pass, but failed with: <ERROR> %s", i+1, actualErr.Error())
		}
		if actualErr == nil && !testCase.shouldPass {
			t.Errorf("Test %d: Expected to fail with <ERROR> \"%s\", but passed instead.", i+1, testCase.expectedErr)
		}
		// Failed as expected, but does it fail for the expected reason.
		if actualErr != nil && !testCase.shouldPass {
			if testCase.expectedErr != actualErr {
				t.Errorf("Test %d: Expected to fail with error \"%s\", but instead failed with error \"%s\" instead.", i+1, testCase.expectedErr, actualErr)
			}
		}
		// Test passes as expected, but the output values
		// are verified for correctness here.
		if actualErr == nil && testCase.shouldPass {
			if testCase.rows != len(actualResult) {
				// End the tests here if the the size doesn't match number of rows.
				t.Fatalf("Test %d: Expected the size of the row of the new matrix to be `%d`, but instead found `%d`", i+1, testCase.rows, len(actualResult))
			}
			// Iterating over each row and validating the size of the column.
			for j, row := range actualResult {
				// If the row check passes, verify the size of each columns.
				if testCase.columns != len(row) {
					t.Errorf("Test %d: Row %d: Expected the size of the column of the new matrix to be `%d`, but instead found `%d`", i+1, j+1, testCase.columns, len(row))
				}
			}
		}
	}
}

// TestMatrixIdentity - validates the method for returning identity matrix of given size.
func TestMatrixIdentityP(t *testing.T) {
	m, err := identityMatrixP(3)
	if err != nil {
		t.Fatal(err)
	}
	str := m.String()
	expect := "[1, 0, 0],\n[0, 1, 0],\n[0, 0, 1]"
	if str != expect {
		t.Fatal(str, "!=", expect)
	}
}

// Tests validate the output of matrix multiplication method.
func TestMatrixMultiply(t *testing.T) {
	m1, err := newMatrixDataP(
		[][]*big.Int{
			{new(big.Int).SetInt64(1), new(big.Int).SetInt64(2)},
			{new(big.Int).SetInt64(3), new(big.Int).SetInt64(4)},
		})
	if err != nil {
		t.Fatal(err)
	}

	m2, err := newMatrixDataP(
		[][]*big.Int{
			{new(big.Int).SetInt64(5), new(big.Int).SetInt64(6)},
			{new(big.Int).SetInt64(7), new(big.Int).SetInt64(8)},
		})
	if err != nil {
		t.Fatal(err)
	}
	actual, err := m1.Multiply(m2, new(big.Int).SetInt64(13))
	if err != nil {
		t.Fatal(err)
	}
	str := actual.String()
	expect := "[6, 9],\n[4, 11]"
	if str != expect {
		t.Fatal(str, "!=", expect)
	}
}

func TestVandermondeP(t *testing.T) {
	m, err := VandermondeP(4, 3, new(big.Int).SetInt64(13))
	if err != nil {
		t.Fatal(err)
	}
	str := m.String()
	expect := "[1, 1, 1],\n[1, 2, 4],\n[1, 3, 9],\n[1, 4, 3]"
	if str != expect {
		t.Fatal(str, "!=", expect)
	}
}
