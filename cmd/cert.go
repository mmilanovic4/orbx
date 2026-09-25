package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var certFile string

func formatCertInfo(cert *x509.Certificate) {
	now := time.Now()
	expires := cert.NotAfter.Format("2006-01-02")

	var expiryStr string
	if now.After(cert.NotAfter) {
		expiryStr = fmt.Sprintf("%s (expired %d days ago) ❌", expires, int(now.Sub(cert.NotAfter).Hours()/24))
	} else {
		daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)
		expiryStr = fmt.Sprintf("%s (%d days left)", expires, daysLeft)
		if daysLeft < 30 {
			expiryStr += " ⚠️"
		}
	}

	fmt.Printf("Subject:    %s\n", cert.Subject.CommonName)
	fmt.Printf("Issuer:     %s\n", cert.Issuer.CommonName)
	fmt.Printf("Valid from: %s\n", cert.NotBefore.Format("2006-01-02"))
	fmt.Printf("Expires:    %s\n", expiryStr)

	if len(cert.DNSNames) > 0 {
		fmt.Printf("SANs:       %s\n", strings.Join(cert.DNSNames, ", "))
	}
}

// splitCertTarget accepts example.com, example.com:8443, [::1]:8443 or a
// URL and returns the host and port to connect to (443 by default).
func splitCertTarget(target string) (string, string) {
	if u, err := url.Parse(target); err == nil && u.Host != "" {
		target = u.Host
	}
	if host, port, err := net.SplitHostPort(target); err == nil {
		return host, port
	}
	return strings.Trim(target, "[]"), "443"
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

		target := args[0]
		host, port := splitCertTarget(target)

		// the handshake skips verification so that expired or self-signed
		// certificates can still be shown, they are verified below instead
		dialer := &net.Dialer{Timeout: 10 * time.Second}
		conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		})
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %w", target, err)
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			return fmt.Errorf("no certificates found for %s", target)
		}

		formatCertInfo(certs[0])

		opts := x509.VerifyOptions{
			DNSName:       host,
			Intermediates: x509.NewCertPool(),
		}
		for _, c := range certs[1:] {
			opts.Intermediates.AddCert(c)
		}
		if _, err := certs[0].Verify(opts); err != nil {
			return fmt.Errorf("certificate is not valid: %w", err)
		}

		return nil
	},
}

func init() {
	certCmd.Flags().StringVarP(&certFile, "file", "f", "", "read certificate from PEM file")
	rootCmd.AddCommand(certCmd)
}
