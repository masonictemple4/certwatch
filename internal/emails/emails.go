package emails

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

var CRLF string = "\r\n"

func SendEmail(targets []string, subject, body string) error {

	// Host
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Source
	from := os.Getenv("EMAIL")
	pw := os.Getenv("EMAIL_PW")

	// Email setup
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	targetStr := strings.Join(targets, ", ")

	toHeader := fmt.Sprintf("To: %s", targetStr+CRLF)
	fromHeader := fmt.Sprintf("From: %s", from+CRLF)
	subjectHeader := fmt.Sprintf("Subject: %s", subject+CRLF)

	emailBody := body + CRLF

	msg := []byte(fromHeader + toHeader + subjectHeader + mime + CRLF + emailBody)

	auth := smtp.PlainAuth("", from, pw, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, targets, msg); err != nil {
		return err
	}

	return nil
}

func InvalidOrExpiredCertEmail(host, body string) error {

	recpStr := os.Getenv("REPORT_RECIPIENTS")

	subStr := fmt.Sprintf("%s - Invalid SSL Certificate [Action Required]", host)

	return SendEmail([]string{recpStr}, subStr, body)
}

func CertExpirationReminder(host string, dayCnt int, body string) error {

	var subStr string
	if dayCnt <= 6 {
		subStr = fmt.Sprintf("%s - SSL Certificate Expires %d days  [Urgent Action Required]", host, dayCnt)
	} else {
		subStr = fmt.Sprintf("%s - SSL Certificate Expires %d days  [Action Required]", host, dayCnt)
	}

	recpStr := os.Getenv("REPORT_RECIPIENTS")

	return SendEmail([]string{recpStr}, subStr, body)
}

func ExpiredCert(host string) error {

	subStr := fmt.Sprintf("%s - SSL Certificate Expired [Urgent Action Required]", host)
	recpStr := os.Getenv("REPORT_RECIPIENTS")

	body := fmt.Sprintf("Greetings!<br />This is an alert to let you know that your SSL certificate for %s has officially expired.<br /><br />To prevent user interruption and unexpected behavior please get the latest Certificate from your Domain provider and upload it to %s", host, host)

	return SendEmail([]string{recpStr}, subStr, body)
}
