package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/newsletter"
)

// subscribeRequest is the newsletter sign-up payload. Consent MUST be true
// (explicit opt-in) — the front-end checkbox is unchecked by default.
type subscribeRequest struct {
	Email   string `json:"email"`
	Lang    string `json:"lang"`
	Consent bool   `json:"consent"`
	Company string `json:"company"` // honeypot — must stay empty
}

// handleSubscribe starts the double opt-in flow: store a pending subscriber and
// e-mail them a confirmation link. Nothing else is ever sent until they confirm.
func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if s.news == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "indisponible"})
		return
	}
	var req subscribeRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<10)).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "requête invalide"})
		return
	}
	// Honeypot: pretend success, do nothing.
	if strings.TrimSpace(req.Company) != "" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "pending"})
		return
	}
	if !req.Consent {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": "Merci de cocher la case de consentement.",
		})
		return
	}
	if !validEmail(req.Email) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": "Merci d’indiquer une adresse e-mail valide.",
		})
		return
	}

	token, already, err := s.news.Subscribe(req.Email, req.Lang)
	if err != nil {
		slog.Error("newsletter subscribe failed", "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "erreur"})
		return
	}
	if !already && token != "" {
		lang := "fr"
		if req.Lang == "en" {
			lang = "en"
		}
		confirmURL := s.cfg.BaseURL + "/newsletter/confirm?token=" + url.QueryEscape(token)
		go func(email, lang, link string) {
			if err := s.sendConfirmationEmail(email, lang, link); err != nil {
				slog.Error("newsletter confirmation email failed", "err", err)
			}
		}(newsletter.NormalizeEmail(req.Email), lang, confirmURL)
	}

	// Always the same response — never leak whether the address already exists.
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "pending",
		"message": "Presque terminé ! Vérifie ta boîte mail et clique sur le lien de confirmation.",
	})
}

// handleNewsletterConfirm completes the double opt-in when the emailed link is
// visited. Idempotent and safe to re-open.
func (s *Server) handleNewsletterConfirm(w http.ResponseWriter, r *http.Request) {
	if s.news == nil {
		http.NotFound(w, r)
		return
	}
	_, lang, ok, err := s.news.Confirm(r.URL.Query().Get("token"))
	if err != nil {
		slog.Error("newsletter confirm failed", "err", err)
	}
	if lang != "en" {
		lang = "fr"
	}
	p := confirmPageFR
	if lang == "en" {
		p = confirmPageEN
	}
	if !ok {
		// Unknown/expired token (or already confirmed) — show a neutral page.
		p = confirmUnknownFR
		if lang == "en" {
			p = confirmUnknownEN
		}
	}
	renderNewsletterPage(w, p, s.cfg.BaseURL)
}

// handleUnsubscribe removes the subscriber (right to erasure). Supports GET
// (link click → HTML page) and POST (RFC 8058 List-Unsubscribe one-click).
func (s *Server) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	if s.news == nil {
		http.NotFound(w, r)
		return
	}
	if _, err := s.news.Unsubscribe(r.URL.Query().Get("token")); err != nil {
		slog.Error("newsletter unsubscribe failed", "err", err)
	}
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusOK) // one-click clients ignore the body
		return
	}
	page := unsubPageFR
	if r.URL.Query().Get("lang") == "en" {
		page = unsubPageEN
	}
	renderNewsletterPage(w, page, s.cfg.BaseURL)
}

// --- e-mail -----------------------------------------------------------------

// smtpConfigured reports whether outbound mail can be sent.
func (s *Server) smtpConfigured() bool {
	return s.cfg.SMTPHost != "" && s.cfg.SMTPUser != ""
}

// sendNewsletterMail sends a plain-text message to a single recipient. extra
// holds already-sanitised header lines (e.g. List-Unsubscribe). No-op if SMTP
// is unconfigured.
func (s *Server) sendNewsletterMail(to, subject, body string, extra []string) error {
	if !s.smtpConfigured() {
		return nil
	}
	c := s.cfg
	from := c.MailFrom
	if from == "" {
		from = c.SMTPUser
	}
	var b strings.Builder
	fmt.Fprintf(&b, "From: Alexandre Iwanesko <%s>\r\n", hdr(from))
	fmt.Fprintf(&b, "To: %s\r\n", hdr(to))
	fmt.Fprintf(&b, "Reply-To: %s\r\n", hdr(c.MailTo))
	fmt.Fprintf(&b, "Subject: %s\r\n", hdr(subject))
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	for _, h := range extra {
		b.WriteString(hdr(h) + "\r\n")
	}
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	b.WriteString(body)

	auth := smtp.PlainAuth("", c.SMTPUser, c.SMTPPass, c.SMTPHost)
	return smtp.SendMail(c.SMTPHost+":"+c.SMTPPort, auth, from, []string{to}, []byte(b.String()))
}

func (s *Server) sendConfirmationEmail(email, lang, confirmURL string) error {
	var subject, body string
	if lang == "en" {
		subject = "Confirm your subscription — Alexandre Iwanesko"
		body = "Hi,\r\n\r\nPlease confirm you want to receive news from iwanesko.ch " +
			"(new articles, weekly tactics, upcoming events) by opening this link:\r\n\r\n" +
			confirmURL + "\r\n\r\nIf you didn’t request this, just ignore this e-mail — " +
			"nothing will be sent without your confirmation.\r\n\r\nAlexandre Iwanesko\r\n"
	} else {
		subject = "Confirme ton inscription — Alexandre Iwanesko"
		body = "Bonjour,\r\n\r\nMerci de confirmer que tu souhaites recevoir les nouveautés " +
			"d’iwanesko.ch (nouveaux articles, tactiques de la semaine, dates de stages) " +
			"en ouvrant ce lien :\r\n\r\n" + confirmURL + "\r\n\r\nSi tu n’es pas à l’origine " +
			"de cette demande, ignore simplement cet e-mail — rien ne sera envoyé sans ta " +
			"confirmation.\r\n\r\nAlexandre Iwanesko\r\n"
	}
	return s.sendNewsletterMail(email, subject, body, nil)
}

// --- confirmation / unsubscribe result pages --------------------------------

type nlPage struct{ Title, Heading, Body string }

var (
	confirmPageFR = nlPage{"Inscription confirmée", "C’est confirmé ✓",
		"Merci ! Tu recevras désormais les nouveautés du site. Tu peux te désabonner à tout moment via le lien en bas de chaque e-mail."}
	confirmPageEN = nlPage{"Subscription confirmed", "You’re in ✓",
		"Thanks! You’ll now receive updates from the site. You can unsubscribe anytime via the link at the bottom of every e-mail."}
	confirmUnknownFR = nlPage{"Lien invalide", "Lien invalide ou expiré",
		"Ce lien de confirmation n’est plus valide. Si tu viens de t’inscrire, utilise le dernier e-mail reçu, ou réinscris-toi."}
	confirmUnknownEN = nlPage{"Invalid link", "Invalid or expired link",
		"This confirmation link is no longer valid. If you just signed up, use the latest e-mail you received, or sign up again."}
	unsubPageFR = nlPage{"Désinscription", "Tu es désabonné·e",
		"Ton adresse a été supprimée de la liste. Tu ne recevras plus d’e-mails. À bientôt !"}
	unsubPageEN = nlPage{"Unsubscribed", "You’re unsubscribed",
		"Your address has been removed from the list. You won’t receive any more e-mails. See you around!"}
)

func renderNewsletterPage(w http.ResponseWriter, p nlPage, baseURL string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Robots-Tag", "noindex")
	_ = nlPageTmpl.Execute(w, map[string]string{
		"Title": p.Title, "Heading": p.Heading, "Body": p.Body, "BaseURL": baseURL,
	})
}

var nlPageTmpl = template.Must(template.New("nl").Parse(`<!doctype html>
<html lang="fr"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="robots" content="noindex"><title>{{.Title}}</title>
<style>
 body{font:16px/1.6 system-ui,sans-serif;color:#1e293b;background:#f8fafc;margin:0;
   display:flex;min-height:100vh;align-items:center;justify-content:center;padding:1.5rem}
 .card{max-width:32rem;background:#fff;border:1px solid #e2e8f0;border-radius:1rem;
   padding:2.5rem;box-shadow:0 10px 30px -12px rgba(15,23,42,.15);text-align:center}
 h1{font-size:1.5rem;margin:0 0 .75rem} p{color:#475569;margin:0 0 1.5rem}
 a{display:inline-block;background:#1e293b;color:#fff;text-decoration:none;
   padding:.7rem 1.4rem;border-radius:.6rem;font-weight:600}
</style></head><body>
<div class="card"><h1>{{.Heading}}</h1><p>{{.Body}}</p>
<a href="{{.BaseURL}}">← iwanesko.ch</a></div>
</body></html>`))
