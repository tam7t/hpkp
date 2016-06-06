package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tam7t/hpkp"
)

func main() {
	cmd := "error"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "example":
		example()
	case "cert":
		cert()
	case "headers":
		headers()
	default:
		fmt.Println("usage: view the source code")
	}
}

func example() {
	s := hpkp.NewMemStorage()
	s.Add("github.com", &hpkp.Header{
		Permanent:  true,
		Sha256Pins: []string{},
	})
	client := &http.Client{}
	client.Transport = &http.Transport{
		DialTLS: hpkp.NewPinDialer(s, true, nil),
	}

	req, err := http.NewRequest("GET", "https://www.github.com", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.StatusCode)
}

func cert() {
	file := os.Args[2]
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	certs, err := x509.ParseCertificates(contents)
	if err != nil {
		log.Fatal(err)
	}

	for i := range certs {
		fmt.Println(hpkp.Fingerprint(certs[i]))
	}
}

func headers() {
	addr := os.Args[2]

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(addr)
	if err != nil {
		log.Fatal(err)
	}

	h := hpkp.ParseHeader(resp)
	j, _ := json.Marshal(h)
	fmt.Println(string(j))
}
