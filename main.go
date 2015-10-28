package main

import (
	"cert-tool/csr"
	"cert-tool/selfsigned"
	"flag"
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Setting up command line parameter flags.
	pSelfSignedOpt := flag.Bool("selfsigned", false, "Self signed certificate generation.")
	pKeySize := flag.Int("size", 2048, "Key size in bits.")
	pSignCert := flag.Bool("sign", false, "Process certificate signing request, and produce a signed cert.")
	pCertificateSignReq := flag.Bool("csr", false, "Create a certificate signing request.")
	pSigningCert := flag.String("signing-cert", "", "Path to certificate to use as a signing cert.")
	pDays := flag.Int("days", 365, "Number of days that cert will be valid.")
	positionalArgs := flag.Args()
	flag.Parse()
	
	// Testing output
	fmt.Println("selfsigned: ", *pSelfSignedOpt)
	fmt.Println("size: ", *pKeySize)
	fmt.Println("sign: ", *pSignCert)
	fmt.Println("csr: ", *pCertificateSignReq)
	fmt.Println("signing-cert: ", *pSigningCert)
	fmt.Println("days: ", *pDays)
	fmt.Println("args: ", positionalArgs)
	
	// Test code
	generateCSR("prvqenam102.namdom002.lab", "prvqenam102.csr")
	generateSelfSignedCert("prvqenam105.namdom002.lab", "prvqenam105.crt.der")
}

func generateCSR(fqdn, filename string) {
	//TODO: Generate a CSR by prompting the user for required data.
	if filename == "" {
		filename = "newcert.csr"
	}
	cert := csr.MakeCertSignReq(fqdn)
	err := ioutil.WriteFile(filename, cert, 0644)
	check(err)
}

func generateSelfSignedCert(fqdn, filename string) {
	if filename == "" {
		filename = "newcert.pem"
	}
	//TODO: Generate a self-signed cert by prompting the user for required data.
	cert := selfsigned.MakeSelfSignedCert(fqdn, false)
	err := ioutil.WriteFile(filename, cert, 0644)
	check(err)
}

func signCert(filename string) {
	//TODO: Process a CSR, and sign the certificate.
}
