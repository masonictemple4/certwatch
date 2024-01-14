package scheduler

import (
	"log/slog"
	"time"

	"github.com/masoncitemple4/certwatch/internal/emails"
	"github.com/masoncitemple4/certwatch/internal/validator"
)

func checkCertsJob(srv *Scheduler, hosts []Host) {
	for _, host := range hosts {
		if host.Hostname == "" {
			continue
		}
		checkCert(srv, host.Hostname, host.Port)
	}
}

func checkCert(srv *Scheduler, hostname, port string) {

	if hostname == "" {
		srv.logger.Error("hostname is empty")
		return
	}

	if port == "" {
		port = "443"
	}

	certResult := validator.Verify(hostname, port)

	previewStr := certResult.PreviewString()

	if !certResult.IsValid() {
		if err := emails.InvalidOrExpiredCertEmail(hostname, previewStr); err != nil {
			srv.logger.Error("failed to send invalid cert email", slog.String("hostname", hostname), slog.String("error", err.Error()))
		}
		return
	}

	now := time.Now()
	y, m, d := now.Date()
	eY, eM, eD := certResult.ExpiresAt.Date()

	if y == eY && m == eM && d == eD {
		daysRemaining := int(certResult.ExpiresAt.Sub(now) / validator.DAY)
		if err := emails.CertExpirationReminder(hostname, daysRemaining, previewStr); err != nil {
			srv.logger.Error("failed to send cert expiration reminder", slog.String("hostname", hostname), slog.String("error", err.Error()))
		}
	}

	srv.logger.Info("cert check complete", slog.String("hostname", hostname), slog.String("cert", previewStr))

}
