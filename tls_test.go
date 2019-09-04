package learning

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

var validityFormat = "2006-01-02"

func TestGetCert(t *testing.T) {
	config := &tls.Config{}
	dialer := &net.Dialer{
		Timeout: time.Second * 5,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", "apps.e2open.com:443", config)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	s := conn.ConnectionState()

	if os.Getenv("PRINT_CERTS") == "" {
		t.Skip("Set the PRINT_CERTS flag to print out the certs.")
	}

	verifiedChains := s.VerifiedChains
	for i, chain := range verifiedChains {
		introspectChain(chain, i+1)
	}
}

func TestGetCertWithPool(t *testing.T) {
	config := &tls.Config{
		RootCAs: getCertPool(t),
	}
	dialer := &net.Dialer{
		Timeout: time.Second * 5,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", "apps.e2open.com:443", config)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	s := conn.ConnectionState()

	if os.Getenv("PRINT_CERTS") == "" {
		t.Skip("Set the PRINT_CERTS flag to print out the certs.")
	}

	verifiedChains := s.VerifiedChains
	for i, chain := range verifiedChains {
		introspectChain(chain, i+1)
	}
}

func getCertPool(t *testing.T) *x509.CertPool {
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM([]byte(goDaddyRoot))
	if !ok {
		t.Fatal("could not parse goDaddyRoot")
	}
	return pool
}

func printCert(cert *x509.Certificate) {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	fmt.Println(string(certPEM))
}

func introspectChain(chain []*x509.Certificate, id int) {
	fmt.Println("Chain:", id)
	for _, cert := range chain {
		fmt.Printf("\t%v sub: %v, validity: %v to %v\n", cert.SerialNumber, cert.Subject,
			cert.NotBefore.Format(validityFormat), cert.NotAfter.Format(validityFormat))
		// printCert(cert)
	}
	fmt.Println("---------------------------------")
}

const goDaddyRoot = `-----BEGIN CERTIFICATE-----
MIIDxTCCAq2gAwIBAgIBADANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMCVVMx
EDAOBgNVBAgTB0FyaXpvbmExEzARBgNVBAcTClNjb3R0c2RhbGUxGjAYBgNVBAoT
EUdvRGFkZHkuY29tLCBJbmMuMTEwLwYDVQQDEyhHbyBEYWRkeSBSb290IENlcnRp
ZmljYXRlIEF1dGhvcml0eSAtIEcyMB4XDTA5MDkwMTAwMDAwMFoXDTM3MTIzMTIz
NTk1OVowgYMxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdBcml6b25hMRMwEQYDVQQH
EwpTY290dHNkYWxlMRowGAYDVQQKExFHb0RhZGR5LmNvbSwgSW5jLjExMC8GA1UE
AxMoR28gRGFkZHkgUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgLSBHMjCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL9xYgjx+lk09xvJGKP3gElY6SKD
E6bFIEMBO4Tx5oVJnyfq9oQbTqC023CYxzIBsQU+B07u9PpPL1kwIuerGVZr4oAH
/PMWdYA5UXvl+TW2dE6pjYIT5LY/qQOD+qK+ihVqf94Lw7YZFAXK6sOoBJQ7Rnwy
DfMAZiLIjWltNowRGLfTshxgtDj6AozO091GB94KPutdfMh8+7ArU6SSYmlRJQVh
GkSBjCypQ5Yj36w6gZoOKcUcqeldHraenjAKOc7xiID7S13MMuyFYkMlNAJWJwGR
tDtwKj9useiciAF9n9T521NtYJ2/LOdYq7hfRvzOxBsDPAnrSTFcaUaz4EcCAwEA
AaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYE
FDqahQcQZyi27/a9BUFuIMGU2g/eMA0GCSqGSIb3DQEBCwUAA4IBAQCZ21151fmX
WWcDYfF+OwYxdS2hII5PZYe096acvNjpL9DbWu7PdIxztDhC2gV7+AJ1uP2lsdeu
9tfeE8tTEH6KRtGX+rcuKxGrkLAngPnon1rpN5+r5N9ss4UXnT3ZJE95kTXWXwTr
gIOrmgIttRD02JDHBHNA7XIloKmf7J6raBKZV8aPEjoJpL1E/QYVN8Gb5DKj7Tjo
2GTzLH4U/ALqn83/B2gX2yKQOC16jdFU8WnjXzPKej17CuPKf1855eJ1usV2GDPO
LPAvTK33sefOT6jEm0pUBsV/fdUID+Ic/n4XuKxe9tQWskMJDE32p2u0mYRlynqI
4uJEvlz36hz1
-----END CERTIFICATE-----`
