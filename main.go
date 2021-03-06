package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CertSignRequest struct {
	Country            string `json:"country"`
	State              string `json:"state"`
	City               string `json:"city"`
	Organization       string `json:"organization"`
	OrganizationalUnit string `json:"organizationalUnit"`
	CommonName         string `json:"commonName"`
	EmailAddress       string `json:"emailAddress"`
	Days               string `json:"days"`
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var certSignRequests []CertSignRequest

const NUM_DAYS_DEFAULT = "730"
const DEFAULT_KEY_SIZE = 2048

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/gencert", createCertBundle).Methods("POST")
	router.HandleFunc("/gencert", allowCORS).Methods("OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func allowCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	createCertBundle(w, r)
	return
}

func createCertBundle(w http.ResponseWriter, r *http.Request) {
	// Create an OpenSSL cert bundle with the given data.
	var certSignRequest CertSignRequest
	_ = json.NewDecoder(r.Body).Decode(&certSignRequest)
	fmt.Printf("[DEBUG]: %s\n", certSignRequest)
	if certSignRequest.Days == "" {
		certSignRequest.Days = NUM_DAYS_DEFAULT
	}
	json.NewEncoder(w).Encode(certSignRequest)

	// Call each of the four functions to create the cert bundle.
	signingConfFile := modifySignOnlyConf(certSignRequest.CommonName)
	generateKeys(certSignRequest.CommonName)
	generateCsr(certSignRequest)
	generateSignedCert(certSignRequest.CommonName, signingConfFile)
	generatePkcs12(certSignRequest.CommonName)
	filename := generateTarball(certSignRequest.CommonName)
	if filename == "" {
		http.Error(w, "Failure to create/retrieve tarball.", 400)
		return
	}
	fmt.Println("Server has built: " + filename)

	// Verify file exists, and open it.
	openfile, err := os.Open(filename)
	defer openfile.Close() // Close file handle once function returns.
	if err != nil {
		// File not found, sending 400
		http.Error(w, "File was not found.", 400)
		return
	}

	// File found, create and send correct headers

	// Get Content-Type of the file
	// Create a buffer to store the header of the file
	fileheader := make([]byte, 512)
	// Copy headers into fileheader buffer
	openfile.Read(fileheader)
	// Get content type of file
	fileContentType := http.DetectContentType(fileheader)

	// Get the file size
	fileStat, _ := openfile.Stat() // Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	// Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	// Send the file
	// Already read 512 bytes from the file, so reset offset back to 0
	openfile.Seek(0, 0)
	io.Copy(w, openfile) // 'Copy' the file to the client

	cleanup(certSignRequest.CommonName, signingConfFile)
	return
}

func generateKeys(commonName string) {
	// Shell out and generate the private key for given commonName.
	cmd := exec.Command("openssl", "genrsa", "-out", commonName+".key", strconv.Itoa(DEFAULT_KEY_SIZE))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	fmt.Printf("Key Generation Result: %q\n", out.String())
}

func generateCsr(csr CertSignRequest) {
	// Shell out and generate the CSR using the given data.
	cmd := exec.Command(
		"openssl",
		"req",
		"-new",
		"-sha256",
		"-batch",
		"-key", csr.CommonName+".key",
		"-subj", "/C="+csr.Country+"/ST="+csr.State+"/L="+csr.City+"/O="+csr.Organization+"/OU="+csr.OrganizationalUnit+"/CN="+csr.CommonName+"/emailAddress="+csr.EmailAddress,
		"-days", csr.Days,
		"-out", csr.CommonName+".csr")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	fmt.Printf("CSR Generation Result: %q\n", out.String())
}

func modifySignOnlyConf(dnsName string) string {
	// Shell out and replace the DNS_NAME entry in sign-only.conf with a valid name.
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	uniqueString := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}(4)
	copyCommand := exec.Command(
		"cp",
		"sign-only.conf",
		"sign-only-"+uniqueString+".conf")
	var copyOut bytes.Buffer
	copyCommand.Stdout = &copyOut
	err := copyCommand.Run()
	check(err)

	cmd := exec.Command(
		"sed",
		"-i",
		"s/DNS_NAME/"+dnsName+"/g",
		"sign-only-"+uniqueString+".conf")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	check(err)
	fmt.Printf("Modifcation of sign-only.conf result: %q\n", out.String())

	return "sign-only-" + uniqueString + ".conf"
}

func generateSignedCert(commonName string, signingConfigFile string) {
	// Shell out and generate the signed cert using the commonName as an identifier.
	cmd := exec.Command(
		"openssl",
		"ca",
		"-extensions",
		"SAN",
		"-config", signingConfigFile,
		"-batch",
		"-notext",
		"-in", commonName+".csr",
		"-out", commonName+".pem.crt",
		"-passin", "pass:novell")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	if err != nil {
		fmt.Println("Signing the CSR has failed!")
		fmt.Println(err)
	}
	fmt.Printf("Cert Signing Result: %q\n", out.String())
}

func generatePkcs12(commonName string) {
	// Shell out and generate the PKCS12 cert store using the commonName as an identifier.
	cmd := exec.Command(
		"openssl",
		"pkcs12",
		"-export",
		"-chain",
		"-CAfile", "rootCA.cert.pem",
		"-inkey", commonName+".key",
		"-in", commonName+".pem.crt",
		"-out", commonName+".p12",
		"-passout", "pass:novell")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	fmt.Printf("PKCS12 Generation Result: %q\n", out.String())
}

func generateTarball(commonName string) string {
	// Shell out and create GZipped tarball of all relevant data.
	tarballName := commonName + "-cert-bundle.tar.gz"
	cmd := exec.Command(
		"tar",
		"cvzf",
		tarballName,
		commonName+".key",
		commonName+".csr",
		commonName+".pem.crt",
		commonName+".p12")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	fmt.Printf("Tarball Creation Result: %q\n", out.String())
	if _, err := os.Stat(tarballName); os.IsNotExist(err) {
		fmt.Println("File was not created!")
		return ""
	}
	return tarballName
}

/*
Clean up the files generated.
Returns true if successfull, false otherwise.
*/
func cleanup(commonName string, signingConfigFile string) bool {
	// TODO: remove copied sign-only.conf file, and all files releated to commonName
	cmd := exec.Command(
		"rm",
		signingConfigFile,
		commonName+".key",
		commonName+".csr",
		commonName+".p12",
		commonName+".pem.crt",
		commonName+"-cert-bundle.tar.gz")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	return false
}
