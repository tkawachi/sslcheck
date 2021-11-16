package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"
)

// Print SSL certificate information
func fetchCertificates(host string) ([]*x509.Certificate, error) {
	//conn, err := tls.Dial("tcp", host+":443", nil)
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", host+":443", nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	state := conn.ConnectionState()
	certs := state.PeerCertificates

	return certs, nil
}

func printIssuer(host string) {
	certs, err := fetchCertificates(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error for host %s: %v\n", host, err)
		return
	}

	if len(certs) == 0 {
		fmt.Fprintf(os.Stderr, "No certificate found for host: %s", host)
		return
	}
	fmt.Printf("%s\t%v\t%s\n", host, certs[0].NotAfter, certs[0].Issuer.CommonName)
}

func removeLastDot(s string) string {
	if len(s) > 0 && s[len(s)-1] == '.' {
		return s[:len(s)-1]
	}
	return s
}

func main() {
	if len(os.Args) > 1 {
		for _, host := range os.Args[1:] {
			printIssuer(removeLastDot(host))
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			printIssuer(removeLastDot(scanner.Text()))
		}
	}
}
