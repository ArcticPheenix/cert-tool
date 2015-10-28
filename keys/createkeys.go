// Package keys creates private/public keypairs with various PKI algorithms.
package keys

import (
	"crypto"
	"crypto/rsa"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

// Creates a RSA private/public keypair. Returns the private key as a byte slice.
func MakeRSAKeypair(name string) []byte {
	// TODO - Implement functionality to create and return a RSA private key.
}

// Creates a DSA private/public keypair. Returns the private key as a byte slice.
func MakeDSAKeypair(name string) []byte {
	// TODO - Implement functionality to create and return a DSA private key.
}

// Creates a ECDSA private/public keypair. Returns the private key as a byte slice.
func MakeECDSAKeypair(name string) []byte {
	// TODO - Implement functionality to create and return a ECDSA private key.
}

// Creates an elliptic curve private/public keypair. Returns the private key as a byte slice.
func MakeEllipticKeypair(name string) []byte {
	// TODO - Implement functionality to create and return an elliptic curve private key.
}