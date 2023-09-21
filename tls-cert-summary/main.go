package main

import (
	"crypto/tls"
	"flag"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	var handler = slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

	var domain, port, sni string
	flag.StringVar(&domain, "domain", "", "Domain to connect to")
	flag.StringVar(&port, "port", "443", "Override port to connect to")
	flag.StringVar(&sni, "sni", "", "Override SNI domain (default matches -domain)")
	flag.Parse()

	if sni == "" {
		sni = domain
	}

	var conf = tls.Config{
		InsecureSkipVerify: true,
		ServerName:         sni,
	}

	conn, err := tls.Dial("tcp", domain+":"+port, &conf)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	var certs = conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		var sans = cert.DNSNames
		sort.Strings(sans)
		slog.Info(
			"certificate_info",
			"issuer", strings.Join(cert.Issuer.Organization, ", "),
			"common_name", cert.Subject.CommonName,
			"subject_alternative_names", sans,
			"start", cert.NotBefore.Format("2006-01-02"),
			"end", cert.NotAfter.Format("2006-01-02"),
			"remaining_days", int(time.Until(cert.NotAfter).Hours()/24),
		)
		os.Exit(0)
	}
}
