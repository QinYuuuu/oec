package main

import (
	"fmt"
	"math/big"
	. "oec/reedsolomonP"
)

func main() {
	// Example usage
	p := big.NewInt(29)
	oec, _ := NewRSGFp(3, 5, p)
	input := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}

	shares, err2 := oec.Encode(input)
	if err2 != nil {
		return
	}
	for _, share := range shares {
		fmt.Println(share)
	}
	shares1 := shares[1:]
	err := oec.Rebuild(shares1, func(share Share) {
		fmt.Printf("Rebuilt share: %d, Data: %s\n", share.Number, share.Data.String())
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
