package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// Print SSL certificate information
func fetchCertificates(host string) ([]*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", host+":443", nil)
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
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}

	if len(certs) == 0 {
		fmt.Fprintf(os.Stderr, "No certificate found for host: %v", host)
		return
	}
	fmt.Printf("%s\t%v\t%s\n", host, certs[0].NotAfter, certs[0].Issuer.CommonName)
}

func main() {
	if len(os.Args) > 1 {
		for _, host := range os.Args[1:] {
			printIssuer(host)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			printIssuer(scanner.Text())
		}
	}

}
