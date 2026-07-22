package server

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"strings"
	"time"
)

// Self-made anti-spam: every form must carry a short-lived, HMAC-signed token
// obtained from GET /api/form-token. A blind bot POSTing straight to the API
// (the vast majority of form spam) has no valid token and is rejected — no
// third-party captcha, no cookie, no tracking.

const formTokenMaxAge = 6 * time.Hour

func newFormKey() []byte {
	k := make([]byte, 32)
	_, _ = rand.Read(k)
	return k
}

func (s *Server) formMAC(msg []byte) []byte {
	m := hmac.New(sha256.New, s.formKey)
	m.Write(msg)
	return m.Sum(nil)[:16]
}

// issueFormToken returns base64url(unix8 || hmac16).
func (s *Server) issueFormToken() string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(time.Now().Unix()))
	return base64.RawURLEncoding.EncodeToString(append(buf, s.formMAC(buf)...))
}

// validFormToken verifies the signature and that the token isn't ancient.
func (s *Server) validFormToken(tok string) bool {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(tok))
	if err != nil || len(raw) != 24 {
		return false
	}
	if !hmac.Equal(s.formMAC(raw[:8]), raw[8:]) {
		return false
	}
	age := time.Since(time.Unix(int64(binary.BigEndian.Uint64(raw[:8])), 0))
	return age >= 0 && age <= formTokenMaxAge
}

func (s *Server) handleFormToken(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]string{"token": s.issueFormToken()})
}

// looksSpammy flags free-text that carries multiple links — a reliable spam
// signal for a coaching contact form (students don't paste URLs).
func looksSpammy(text string) bool {
	return strings.Count(strings.ToLower(text), "http") >= 2
}
