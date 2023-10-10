package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type certificateInfo struct {
	Issuer                   string   `json:"issuer,omitempty"`
	CommonName               string   `json:"common_name,omitempty"`
	SubjectAlternateiveNames []string `json:"subject_alternative_names,omitempty"`
	Start                    string   `json:"start,omitempty"`
	End                      string   `json:"end,omitempty"`
	RemainingDays            int      `json:"remaining_days,omitempty"`
}

const dateFormat = "2006-01-02"

func main() {
	var domain, port, sni string

	flag.StringVar(&domain, "domain", "", "Domain to connect to")
	flag.StringVar(&port, "port", "443", "Override port to connect to")
	flag.StringVar(&sni, "sni", "", "Override SNI domain (default matches -domain)")
	flag.Parse()

	if domain == "" {
		log.Fatal("please provide a domain using the -domain flag")
	}

	if sni == "" {
		sni = domain
	}

	var conf = tls.Config{
		InsecureSkipVerify: true,
		ServerName:         sni,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", domain, port), &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	var sans = cert.DNSNames
	sort.Strings(sans)

	certInfo := certificateInfo{
		Issuer:                   strings.Join(cert.Issuer.Organization, ", "),
		CommonName:               cert.Subject.CommonName,
		SubjectAlternateiveNames: sans,
		Start:                    cert.NotBefore.Format(dateFormat),
		End:                      cert.NotAfter.Format(dateFormat),
		RemainingDays:            int(time.Until(cert.NotAfter).Hours() / 24),
	}

	if err := json.NewEncoder(os.Stdout).Encode(certInfo); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
