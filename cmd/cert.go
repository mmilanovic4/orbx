package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var certFile string

func formatCertInfo(cert *x509.Certificate) {
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)

	expiryStr := fmt.Sprintf("%s (%d days left)", cert.NotAfter.Format("2006-01-02"), daysLeft)
	if daysLeft < 30 {
		expiryStr += " ⚠️"
	}
	if daysLeft < 0 {
		expiryStr = fmt.Sprintf("%s (expired %d days ago) ❌", cert.NotAfter.Format("2006-01-02"), -daysLeft)
	}

	fmt.Printf("Subject:    %s\n", cert.Subject.CommonName)
	fmt.Printf("Issuer:     %s\n", cert.Issuer.CommonName)
	fmt.Printf("Valid from: %s\n", cert.NotBefore.Format("2006-01-02"))
	fmt.Printf("Expires:    %s\n", expiryStr)

	if len(cert.DNSNames) > 0 {
		fmt.Printf("SANs:       %s\n", strings.Join(cert.DNSNames, ", "))
	}
}

var certCmd = &cobra.Command{
	Use:     "cert [domain]",
	Short:   "Show TLS certificate info for a domain or file",
	GroupID: "network",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if certFile != "" {
			data, err := os.ReadFile(certFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			block, _ := pem.Decode(data)
			if block == nil {
				return fmt.Errorf("failed to decode PEM block")
			}

			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse certificate: %w", err)
			}

			formatCertInfo(cert)
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("domain or --file required")
		}

		domain := args[0]
		host := domain + ":443"

		conn, err := tls.Dial("tcp", host, &tls.Config{
			InsecureSkipVerify: false,
		})
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %w", domain, err)
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			return fmt.Errorf("no certificates found for %s", domain)
		}

		formatCertInfo(certs[0])
		return nil
	},
}

func init() {
	certCmd.Flags().StringVarP(&certFile, "file", "f", "", "read certificate from PEM file")
	rootCmd.AddCommand(certCmd)
}
