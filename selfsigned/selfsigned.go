// Package selfsigned creates self-signed certificates.
package selfsigned

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
)

// Creates a self signed certificate for the specified FQDN. Defaults to SHA256WithRSA.
func MakeSelfSignedCert(fqdn string) []byte {
	// TODO - Implement functionality that creates and returns a self-signed certificate.
}