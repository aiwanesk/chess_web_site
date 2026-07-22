package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/booking"
)

// Lesson slots run 17:30–20:00, in 30-minute steps.
const (
	dayStartMin  = 17*60 + 30 // 17:30
	dayEndMin    = 20 * 60    // 20:00
	slotStep     = 30
	minLessonMin = 60 // a lesson lasts at least one hour
)

type bookingRequest struct {
	Date    string `json:"date"`  // YYYY-MM-DD
	Start   string `json:"start"` // "HH:MM"
	End     string `json:"end"`   // "HH:MM"
	Name    string `json:"name"`
	Email   string `json:"email"`
	Lang    string `json:"lang"`
	Token   string `json:"token"`   // anti-spam form token
	Company string `json:"company"` // honeypot
}

// handleBooking validates a lesson reservation, stores it (rejecting overlaps),
// computes the amount due and e-mails both the visitor and Alexandre.
func (s *Server) handleBooking(w http.ResponseWriter, r *http.Request) {
	if s.bookings == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "indisponible"})
		return
	}
	var req bookingRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<10)).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "requête invalide"})
		return
	}
	if strings.TrimSpace(req.Company) != "" { // honeypot
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	if !s.validFormToken(req.Token) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Session expirée, merci de réessayer.", "code": "token"})
		return
	}

	if strings.TrimSpace(req.Name) == "" || !validEmail(req.Email) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Nom et e-mail valides requis."})
		return
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Date invalide."})
		return
	}
	if req.Date < s.bookingFloor() {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Cette date n’est pas encore ouverte à la réservation."})
		return
	}
	start, ok1 := parseHHMM(req.Start)
	end, ok2 := parseHHMM(req.End)
	if !ok1 || !ok2 || start%slotStep != 0 || end%slotStep != 0 ||
		start < dayStartMin || end > dayEndMin || end <= start {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": "Créneau invalide (17h30–20h00, par tranches de 30 min).",
		})
		return
	}
	if end-start < minLessonMin {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Un cours dure au minimum 1 heure."})
		return
	}

	price := s.cfg.HourlyRate * (end - start) / 60

	b := booking.Booking{
		Date: req.Date, StartMin: start, EndMin: end,
		Name: strings.TrimSpace(req.Name), Email: strings.TrimSpace(req.Email), Price: price,
	}
	if err := s.bookings.Create(b); err != nil {
		if errors.Is(err, booking.ErrConflict) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Ce créneau est déjà réservé, choisis-en un autre."})
			return
		}
		slog.Error("booking create failed", "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "erreur"})
		return
	}

	go func(b booking.Booking, lang string) {
		if err := s.sendBookingEmails(b, lang); err != nil {
			slog.Error("booking email failed", "err", err)
		}
	}(b, req.Lang)

	slog.Info("lesson booked", "date", b.Date, "start", b.StartMin, "end", b.EndMin, "price", price)
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"price":   price,
		"message": "Cours réservé ! Tu vas recevoir un e-mail de confirmation.",
	})
}

// bookingFloor is the earliest bookable date (later of today and the configured
// opening date).
func (s *Server) bookingFloor() string {
	today := time.Now().Format("2006-01-02")
	if s.cfg.BookingMinDate > today {
		return s.cfg.BookingMinDate
	}
	return today
}

// handleAvailability returns the already-booked slots for a day (occupied time
// ranges only — never names or e-mails), so the form can grey them out.
func (s *Server) handleAvailability(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if s.bookings == nil {
		writeJSON(w, http.StatusOK, map[string]any{"taken": []any{}})
		return
	}
	date := r.URL.Query().Get("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "date invalide"})
		return
	}
	slots, err := s.bookings.TakenOn(date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "erreur"})
		return
	}
	taken := make([]map[string]string, 0, len(slots))
	for _, sl := range slots {
		taken = append(taken, map[string]string{"start": minToHHMM(sl.StartMin), "end": minToHHMM(sl.EndMin)})
	}
	writeJSON(w, http.StatusOK, map[string]any{"taken": taken})
}

// handleBookingConfig exposes the constraints the booking form needs.
func (s *Server) handleBookingConfig(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{
		"minDate":    s.bookingFloor(),
		"minMinutes": minLessonMin,
		"hourlyRate": s.cfg.HourlyRate,
	})
}

// sendBookingEmails confirms to the visitor and notifies Alexandre.
func (s *Server) sendBookingEmails(b booking.Booking, lang string) error {
	c := s.cfg
	if c.SMTPHost == "" || c.SMTPUser == "" {
		return nil
	}
	from := c.MailFrom
	if from == "" {
		from = c.SMTPUser
	}
	when := fmt.Sprintf("%s, %s–%s", b.Date, minToHHMM(b.StartMin), minToHHMM(b.EndMin))
	dur := durationLabel(b.EndMin-b.StartMin, lang)

	// 1) Confirmation to the visitor (Reply-To = Alexandre).
	var subj, body string
	if lang == "en" {
		subj = "Your chess lesson is booked — " + b.Date
		body = fmt.Sprintf("Hi %s,\r\n\r\nYour lesson is confirmed:\r\n\r\n"+
			"Date & time: %s\r\nDuration: %s\r\nAmount due: %d CHF\r\n\r\n"+
			"I'll be in touch if anything's needed. See you at the board!\r\nAlexandre Iwanesko\r\n",
			b.Name, when, dur, b.Price)
	} else {
		subj = "Ton cours d’échecs est réservé — " + b.Date
		body = fmt.Sprintf("Bonjour %s,\r\n\r\nTon cours est confirmé :\r\n\r\n"+
			"Date et horaire : %s\r\nDurée : %s\r\nMontant à régler : %d CHF\r\n\r\n"+
			"Je te recontacte si besoin. À très vite sur l’échiquier !\r\nAlexandre Iwanesko\r\n",
			b.Name, when, dur, b.Price)
	}
	if err := sendPlain(c, from, b.Email, hdr(c.MailTo), subj, body); err != nil {
		return err
	}

	// 2) Notification to Alexandre (Reply-To = the visitor).
	adminBody := fmt.Sprintf("Nouvelle réservation de cours.\r\n\r\n"+
		"Nom : %s\r\nE-mail : %s\r\nDate/horaire : %s\r\nDurée : %s\r\nMontant : %d CHF\r\n",
		b.Name, b.Email, when, durationLabel(b.EndMin-b.StartMin, "fr"), b.Price)
	return sendPlain(c, from, c.MailTo, hdr(b.Email), "Nouvelle réservation de cours — "+b.Date, adminBody)
}

// sendPlain sends a UTF-8 text/plain message (headers sanitised).
func sendPlain(c Config, from, to, replyTo, subject, body string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "From: Alexandre Iwanesko <%s>\r\n", hdr(from))
	fmt.Fprintf(&b, "To: %s\r\n", hdr(to))
	if replyTo != "" {
		fmt.Fprintf(&b, "Reply-To: %s\r\n", replyTo)
	}
	fmt.Fprintf(&b, "Subject: %s\r\n", hdr(subject))
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	b.WriteString(body)
	auth := smtp.PlainAuth("", c.SMTPUser, c.SMTPPass, c.SMTPHost)
	return smtp.SendMail(c.SMTPHost+":"+c.SMTPPort, auth, from, []string{to}, []byte(b.String()))
}

func parseHHMM(v string) (int, bool) {
	parts := strings.SplitN(strings.TrimSpace(v), ":", 2)
	if len(parts) != 2 {
		return 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, false
	}
	return h*60 + m, true
}

func minToHHMM(m int) string {
	return fmt.Sprintf("%02d:%02d", m/60, m%60)
}

func durationLabel(mins int, lang string) string {
	h, m := mins/60, mins%60
	if lang == "en" {
		if m == 0 {
			return fmt.Sprintf("%dh", h)
		}
		return fmt.Sprintf("%dh%02d", h, m)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh%02d", h, m)
}
