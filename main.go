package main

import (
	"fmt"
	"math/big"
	. "oec/reedsolomonP"
)

func main() {
	// Example usage
	p := big.NewInt(29)
	n := 7
	k := 3
	oec, _ := NewRSGFp(k, n, p)
	input := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(4)}

	shares, err2 := oec.Encode(input)
	if err2 != nil {
		return
	}
	shares[3] = Share{
		Number: 3,
		Data:   big.NewInt(1),
	}

	shares[4] = Share{
		Number: 4,
		Data:   big.NewInt(1),
	}
	for _, share := range shares {
		fmt.Println(share)
	}
	shares2, err := oec.Correct(shares)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(shares2)
}
