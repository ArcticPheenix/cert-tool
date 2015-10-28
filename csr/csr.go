// Package csr creates Certificate Signing Requests
package csr

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
)

// Creates a certificate signing request for the specified FQDN. Defaults to SHA256WithRSA.
func MakeCertSignReq(fqdn string) []byte {
	// Populate the subject data.
	subject := pkix.Name{
		Country:            []string{"US"},
		Organization:       []string{"Microfocus International"},
		OrganizationalUnit: []string{"Identity and Access Management"},
		Locality:           []string{"Provo"},
		Province:           []string{"Utah"},
		CommonName:         fqdn,
	}

	// Populate the cert request template.
	certificateRequest := &x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: 4,
		DNSNames:           []string{fqdn},
		EmailAddresses:     []string{"christopher.kelly@microfocus.com"},
	}

	// Create a private key.
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
	}

	// Create signing request and return byte slice if no errors were encountered.
	certSignReq, err := x509.CreateCertificateRequest(rand.Reader, certificateRequest, privatekey)
	if err != nil {
		fmt.Println(err)
	}
	return certSignReq
}
