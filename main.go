package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/sjhitchner/bloomd/bloom"
)

func main() {

	bloom := bloom.NewSimpleBloomFilter(2)

	for i := 0; i < 25; i++ {
		str := uuid.New()
		bloom.Add(str)
		if bloom.Test(str) {
			fmt.Println("exists!")
		}
	}
	bloom.Add("steve")

	fmt.Printf("False Positives: %f\n", bloom.FalsePositives())
}

/*

0  3 4  7
0000 0000

1 >> 2


*/
