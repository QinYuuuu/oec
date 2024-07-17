package main

import (
	"fmt"
	"oec"
)

func main() {
	f := 1
	n := 4
	coder, err := oec.NewFEC(f+1, n)
	if err != nil {
		panic(err)
	}
	message := []byte{0x01, 0x02}
	shares := make([]oec.Share, n)
	output := func(s oec.Share) {
		// we need to make a copy of the data. The share data
		// memory gets reused when we return.
		shares[s.Number] = s.DeepCopy()
	}
	err = coder.Encode(message, output)
	if err != nil {
		panic(err)
	}
	fmt.Println("shares:", shares)
}
