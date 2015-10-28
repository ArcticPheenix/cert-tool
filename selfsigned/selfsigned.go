// Package selfsigned creates self-signed certificates.
package selfsigned

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

// Creates a self signed certificate for the specified FQDN. Defaults to SHA256WithRSA.
func MakeSelfSignedCert(fqdn string, isCA bool) []byte {
	// TODO - Implement functionality that creates and returns a self-signed certificate.
	randMax := big.NewInt(int64(1000))
	serialNumber, err := rand.Int(rand.Reader, randMax)
	if err != nil {
		fmt.Println(err)
	}
	
	isCA = false
	
	//DEBUG
	fmt.Printf("Value of serialNumber: %v", serialNumber)
	timeVar := time.Now()
	subject := pkix.Name{
		Country:            []string{"US"},
		Organization:       []string{"Microfocus International"},
		OrganizationalUnit: []string{"Identity and Access Management"},
		Locality:           []string{"Provo"},
		Province:           []string{"Utah"},
		CommonName:         fqdn,
	}
	
	certTemplate := &x509.Certificate{
		Subject:            subject,
		NotBefore:          time.Now(),
		NotAfter:           timeVar.AddDate(0,0,365),  // Add 365 days to timeVar timestamp.
		DNSNames:           []string{fqdn},
		EmailAddresses:     []string{"christopher.kelly@microfocus.com"},
		SerialNumber:       serialNumber,
		SignatureAlgorithm: 4,
	}
	
	// Create a private RSA key.
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
	}
	
	// Compute public key from privateKey.
	publickey := privatekey.Public()
	
	// Create certificate and return byte slice if no errors were encountered.
	certificate, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, publickey, privatekey)
	if err != nil {
		fmt.Println(err)
	}
	return certificate
}