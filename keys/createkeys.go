// Package keys creates private/public keypairs with various PKI algorithms.
package keys

import (
	"crypto/dsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// Creates a RSA private/public keypair. Returns the private key as a byte slice.
func MakeRSAKeypair(name string, keysize int) {
	// TODO - Implement functionality to create and return a RSA private key.
	keypair, err := rsa.GenerateKey(rand.Reader, keysize)
	check(err)
}

// Creates a DSA private/public keypair. Returns the private key as a byte slice.
func MakeDSAKeypair(name string, keysize int) {
	// TODO - Implement functionality to create and return a DSA private key.
	keypair, err := dsa.GenerateKey(rand.Reader, keysize)
	check(err)
	fmt.Println(keypair)
}

// Creates a ECDSA private/public keypair. Returns the private key as a byte slice.
func MakeECDSAKeypair(name string, keysize int) {
	// TODO - Implement functionality to create and return a ECDSA private key.
	keypair, err := ecsda.GenerateKey(rand.Reader, keysize)
	check(err)
	fmt.Println(keypair)
}

// Creates an elliptic curve private/public keypair. Returns the private key as a byte slice.
func MakeEllipticKeypair(name string, keysize int) {
	// TODO - Implement functionality to create and return an elliptic curve private key.
	keypair, err := elliptic.GenerateKey(rand.Reader, keysize)
	check(err)
	fmt.Println(keypair)
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
