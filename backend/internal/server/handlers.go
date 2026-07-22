package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// contactRequest is the payload from the frontend contact form.
type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Level   string `json:"level"` // e.g. Elo range or self-assessment
	Goal    string `json:"goal"`  // what the person wants to achieve
	Message string `json:"message"`
	Lang    string `json:"lang"` // "fr" | "en" — for the acknowledgement e-mail
	// Honeypot must stay empty; bots tend to fill every field.
	Honeypot string `json:"company"`
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	var req contactRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "requête invalide"})
		return
	}

	if strings.TrimSpace(req.Honeypot) != "" {
		// Silently accept to avoid teaching bots the honeypot exists.
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	if strings.TrimSpace(req.Name) == "" || !validEmail(req.Email) || strings.TrimSpace(req.Message) == "" {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": "Merci de renseigner votre nom, un e-mail valide et un message.",
		})
		return
	}

	if s.cfg.ContactWebhookURL != "" {
		if err := s.forwardContact(req); err != nil {
			slog.Error("contact forward failed", "err", err)
			// Don't leak infra errors; the lead is logged below regardless.
		}
	}
	// Send the email to Alexandre and the acknowledgement to the visitor
	// asynchronously so the visitor gets a fast response; the lead is logged
	// synchronously so nothing is lost on failure.
	go func(req contactRequest) {
		if err := s.sendContactEmail(req); err != nil {
			slog.Error("contact email failed", "err", err)
		}
		if err := s.sendContactAck(req); err != nil {
			slog.Error("contact acknowledgement failed", "err", err)
		}
	}(req)
	slog.Info("contact lead", "name", req.Name, "email", req.Email, "level", req.Level)

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Merci ! Votre message a bien été envoyé, je vous réponds rapidement.",
	})
}

func (s *Server) forwardContact(req contactRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.ContactWebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// validEmail is a deliberately loose check: presence of a single @ with text
// on both sides and a dot in the domain. Real validation is delivery.
func validEmail(s string) bool {
	s = strings.TrimSpace(s)
	at := strings.IndexByte(s, '@')
	if at <= 0 || at == len(s)-1 {
		return false
	}
	return strings.Contains(s[at+1:], ".")
}
