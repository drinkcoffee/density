package density

import (
    "crypto/rand"
    "math/big"
)


// CreateRandomIdentifier generates a random id in the range 0 to 256**wordSize - 1
// Typically this would be:
// - generate random number (the private key)
// - derive the public key / scalar multiply using the private key as the scalar
// - message digest the public key
func CreateRandomIdentifier(max *big.Int) (*big.Int) {
	n, _ := rand.Int(rand.Reader, max)
	return n
}

