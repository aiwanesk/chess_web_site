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

// sendContactAck sends an acknowledgement e-mail to the visitor confirming that
// their message was received. Reply-To is Alexandre's inbox, so a reply reaches
// him directly. No-op if SMTP isn't configured.
func (s *Server) sendContactAck(req contactRequest) error {
	c := s.cfg
	if c.SMTPHost == "" || c.SMTPUser == "" {
		return nil
	}
	from := c.MailFrom
	if from == "" {
		from = c.SMTPUser
	}
	to := hdr(req.Email)
	if to == "" {
		return nil
	}

	var subject, body string
	if req.Lang == "en" {
		subject = "We got your message — Alexandre Iwanesko"
		body = fmt.Sprintf("Hi %s,\r\n\r\nThanks for reaching out! I’ve received your message "+
			"and I’ll get back to you shortly to discuss how we can work together.\r\n\r\n"+
			"For reference, here is what you sent:\r\n\r\n%s\r\n\r\n"+
			"Talk soon,\r\nAlexandre Iwanesko\r\nFIDE Master · Chess coach\r\n",
			req.Name, strings.TrimSpace(req.Message))
	} else {
		subject = "Votre message est bien reçu — Alexandre Iwanesko"
		body = fmt.Sprintf("Bonjour %s,\r\n\r\nMerci pour votre message ! Je l’ai bien reçu et "+
			"je reviens vers vous rapidement pour échanger sur vos objectifs.\r\n\r\n"+
			"Pour rappel, voici ce que vous m’avez envoyé :\r\n\r\n%s\r\n\r\n"+
			"À très vite,\r\nAlexandre Iwanesko\r\nMaître FIDE · Coach d’échecs\r\n",
			req.Name, strings.TrimSpace(req.Message))
	}

	var b strings.Builder
	fmt.Fprintf(&b, "From: Alexandre Iwanesko <%s>\r\n", hdr(from))
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Reply-To: %s\r\n", hdr(c.MailTo))
	fmt.Fprintf(&b, "Subject: %s\r\n", hdr(subject))
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	b.WriteString(body)

	auth := smtp.PlainAuth("", c.SMTPUser, c.SMTPPass, c.SMTPHost)
	return smtp.SendMail(c.SMTPHost+":"+c.SMTPPort, auth, from, []string{req.Email}, []byte(b.String()))
}

// hdr strips CR/LF from a value used in an email header (anti header-injection).
func hdr(v string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(strings.TrimSpace(v))
}
