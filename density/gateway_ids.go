package main

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "sort"
)

const WORD_SIZE = 32  // 32 bytes
var max = new(big.Int)



type GatewayIds struct {
	ids[] big.Int
}


// Create a new instance
func newGatewayIds() *GatewayIds {
	max.Exp(big.NewInt(2), big.NewInt(8 * WORD_SIZE), nil).Sub(max, big.NewInt(1))
    fmt.Printf("Max: 0x%x\n", max)

	var gatewayIds = GatewayIds{}
	gatewayIds.ids = []big.Int{ } 
	var _ sort.Interface = &gatewayIds  // Enforce interface compliance
	return &gatewayIds
}

func (this *GatewayIds) AddRandom() {
	n, _ := rand.Int(rand.Reader, max)
	this.ids = append(this.ids, *n)
}

// Len is the number of elements in the collection.
func (this *GatewayIds) Len() int {
	return len(this.ids)
}


// Less reports whether the element with
// index i should sort before the element with index j.
func (this *GatewayIds)  Less(i, j int) bool {
	return true
}
	
// Swap swaps the elements with indexes i and j.
func (this *GatewayIds) Swap(i, j int) {

}