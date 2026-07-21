package server

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// sendContactEmail sends the contact submission to MailTo over SMTP (STARTTLS on
// port 587). No-op if SMTP isn't configured — the lead is still logged/forwarded.
func (s *Server) sendContactEmail(req contactRequest) error {
	c := s.cfg
	if c.SMTPHost == "" || c.SMTPUser == "" || c.MailTo == "" {
		return nil
	}
	from := c.MailFrom
	if from == "" {
		from = c.SMTPUser
	}

	// Header values are sanitised (no CR/LF) to prevent header injection.
	subject := hdr("Nouveau message du site — " + req.Name)
	replyTo := hdr(req.Email)

	var b strings.Builder
	fmt.Fprintf(&b, "From: iwanesko.ch <%s>\r\n", hdr(from))
	fmt.Fprintf(&b, "To: %s\r\n", hdr(c.MailTo))
	if replyTo != "" {
		fmt.Fprintf(&b, "Reply-To: %s\r\n", replyTo)
	}
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	// Body (free-form text is safe here — not header context).
	fmt.Fprintf(&b, "Nouveau message via le formulaire de contact.\r\n\r\n")
	fmt.Fprintf(&b, "Nom : %s\r\n", req.Name)
	fmt.Fprintf(&b, "E-mail : %s\r\n", req.Email)
	if strings.TrimSpace(req.Level) != "" {
		fmt.Fprintf(&b, "Niveau : %s\r\n", req.Level)
	}
	fmt.Fprintf(&b, "\r\nMessage :\r\n%s\r\n", req.Message)

	auth := smtp.PlainAuth("", c.SMTPUser, c.SMTPPass, c.SMTPHost)
	addr := c.SMTPHost + ":" + c.SMTPPort
	return smtp.SendMail(addr, auth, from, []string{c.MailTo}, []byte(b.String()))
}

// hdr strips CR/LF from a value used in an email header (anti header-injection).
func hdr(v string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(strings.TrimSpace(v))
}
