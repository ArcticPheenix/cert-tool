package main

import (
	"flag"
	"fmt"
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
}
