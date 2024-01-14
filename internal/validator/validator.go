package validator

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
	"time"
)

var (
	DAY                     = 24 * time.Hour
	WARNING_DURATION_LEGEND = []int{14, 9, 6, 3, 2, 1, 0}
)

type Result struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`

	// Contains all of the necessary information about the connection,
	// and the certificates.
	Conn tls.ConnectionState `json:"conn"`

	// Convenience fields
	IsExpired     bool      `json:"is_expired"`
	ValidHostname bool      `json:"valid_hostname"`
	ExpiresAt     time.Time `json:"expires_at"`

	// errors encountered during validation.
	// might eventually separate app errors and validation errors.
	Errors []error `json:"errors"`
}

func Verify(hostname, port string) *Result {
	result := Result{
		Hostname:      hostname,
		Port:          port,
		ValidHostname: true,
	}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), nil)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return &result
	}

	if err := conn.VerifyHostname(hostname); err != nil {
		result.Errors = append(result.Errors, err)
		result.ValidHostname = false
		return &result
	}

	result.Conn = conn.ConnectionState()

	expiry := result.Conn.PeerCertificates[0].NotAfter

	result.ExpiresAt = expiry

	return &result
}

func (r *Result) IsValid() bool {
	return r.ValidHostname && !r.IsExpired
}

func (r *Result) PreviewString() string {
	var strBldr strings.Builder

	domains := make([]string, 0)
	crtEmails := make([]string, 0)
	var issuer string
	var subject string
	if len(r.Conn.PeerCertificates) > 0 {
		domains = r.Conn.PeerCertificates[0].DNSNames
		crtEmails = r.Conn.PeerCertificates[0].EmailAddresses
		subject = r.Conn.PeerCertificates[0].Subject.String()
		issuer = r.Conn.PeerCertificates[0].Issuer.String()
	}

	fmt.Fprintf(&strBldr, "Host: %s\n", r.Hostname)
	fmt.Fprintf(&strBldr, "Issuer: %s\nExpiry: %v\n", issuer, r.ExpiresAt.Format(time.RFC850))
	fmt.Fprintf(&strBldr, "Subject: %s\n", subject)

	strBldr.WriteString("Domains: \n")
	for _, domain := range domains {
		strBldr.WriteString("- ")
		strBldr.WriteString(domain)
		strBldr.WriteString("\n")
	}

	strBldr.WriteString("Emails: \n")
	for _, email := range crtEmails {
		strBldr.WriteString("- ")
		strBldr.WriteString(email)
		strBldr.WriteString("\n")
	}

	strBldr.WriteString("Errors: \n")
	for _, err := range r.Errors {
		strBldr.WriteString("-")
		strBldr.WriteString(err.Error())
		strBldr.WriteString("\n")
	}

	return strBldr.String()

}

func (r *Result) PeerCertificates() []*x509.Certificate {
	return r.Conn.PeerCertificates
}

func (r *Result) VerifiedChaines() [][]*x509.Certificate {
	return r.Conn.VerifiedChains
}

func (r *Result) RawPeerCertificatesString() string {
	var strBldr strings.Builder

	for _, cert := range r.PeerCertificates() {
		strBldr.WriteString(string(cert.Raw))
	}

	return strBldr.String()
}

func (r *Result) RawVerifiedChainsString() string {
	var strBldr strings.Builder

	for _, chain := range r.VerifiedChaines() {
		for _, cert := range chain {
			// TODO: This loop might be off, will need to inspect.
			// From what I understand PeerCertificates[0] is the leaf cert.
			strBldr.WriteString(string(cert.Raw))
		}
	}

	return strBldr.String()
}

func (r *Result) NextWarningDate() time.Time {
	now := time.Now()
	daysTillExpiry := int(r.ExpiresAt.Sub(now) / DAY)

	for _, days := range WARNING_DURATION_LEGEND {
		if daysTillExpiry >= days {
			return r.ExpiresAt.Add(-time.Duration(days) * DAY)
		}
	}

	return r.ExpiresAt
}
