package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		panic(e)
	}
}

var certSignRequests []CertSignRequest

const NUM_DAYS_DEFAULT = "730"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/gencert", CreateCertBundle).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func CreateCertBundle(w http.ResponseWriter, r *http.Request) {
	//TODO: Create an OpenSSL cert bundle with the given data.
	params := mux.Vars(r)
	fmt.Println(r)
	var certSignRequest CertSignRequest
	_ = json.NewDecoder(r.Body).Decode(&certSignRequest)
	certSignRequest.Country = params["country"]
	certSignRequest.State = params["state"]
	certSignRequest.City = params["city"]
	certSignRequest.Organization = params["organization"]
	certSignRequest.OrganizationalUnit = params["organizationalUnit"]
	certSignRequest.CommonName = params["commonName"]
	certSignRequest.EmailAddress = params["emailAddress"]
	if params["days"] != "" {
		certSignRequest.Days = params["days"]
	} else {
		certSignRequest.Days = "730"
	}
	json.NewEncoder(w).Encode(certSignRequest)
}
