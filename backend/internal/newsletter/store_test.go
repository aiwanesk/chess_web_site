package newsletter

import (
	"path/filepath"
	"testing"
)

func open(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "nl.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestDoubleOptInFlow(t *testing.T) {
	s := open(t)

	ct, already, err := s.Subscribe("  Alice@Example.com ", "fr")
	if err != nil || already || ct == "" {
		t.Fatalf("subscribe: ct=%q already=%v err=%v", ct, already, err)
	}

	// Pending subscriber must NOT be a recipient yet (explicit consent pending).
	if rs, _ := s.ConfirmedRecipients(); len(rs) != 0 {
		t.Fatalf("pending should not receive mail, got %d recipients", len(rs))
	}

	email, lang, ok, err := s.Confirm(ct)
	if err != nil || !ok {
		t.Fatalf("confirm: ok=%v err=%v", ok, err)
	}
	if email != "alice@example.com" {
		t.Fatalf("email not normalised: %q", email)
	}
	if lang != "fr" {
		t.Fatalf("lang=%q", lang)
	}

	rs, _ := s.ConfirmedRecipients()
	if len(rs) != 1 || rs[0].Email != "alice@example.com" {
		t.Fatalf("recipients after confirm = %+v", rs)
	}

	// Re-using the confirm token must not re-confirm.
	if _, _, ok, _ := s.Confirm(ct); ok {
		t.Fatal("second confirm should be a no-op")
	}

	// Unsubscribe deletes the row (right to erasure).
	removed, err := s.Unsubscribe(rs[0].UnsubToken)
	if err != nil || !removed {
		t.Fatalf("unsubscribe: removed=%v err=%v", removed, err)
	}
	if rs, _ := s.ConfirmedRecipients(); len(rs) != 0 {
		t.Fatalf("after unsubscribe recipients = %d", len(rs))
	}
}

func TestSubscribeAlreadyConfirmedIsIdempotent(t *testing.T) {
	s := open(t)
	ct, _, _ := s.Subscribe("bob@example.com", "en")
	if _, _, ok, _ := s.Confirm(ct); !ok {
		t.Fatal("confirm failed")
	}
	ct2, already, err := s.Subscribe("bob@example.com", "en")
	if err != nil {
		t.Fatal(err)
	}
	if !already || ct2 != "" {
		t.Fatalf("re-subscribe of confirmed: already=%v token=%q (want true/empty)", already, ct2)
	}
}

func TestSeedAndNotified(t *testing.T) {
	s := open(t)
	if seeded, _ := s.Seeded(); seeded {
		t.Fatal("fresh DB must not be seeded")
	}
	if err := s.MarkNotified("blog:a", "tactics:01-01-25"); err != nil {
		t.Fatal(err)
	}
	if err := s.MarkSeeded(); err != nil {
		t.Fatal(err)
	}
	if seeded, _ := s.Seeded(); !seeded {
		t.Fatal("should be seeded after MarkSeeded")
	}
	if done, _ := s.IsNotified("blog:a"); !done {
		t.Fatal("blog:a should be marked notified")
	}
	if done, _ := s.IsNotified("blog:brand-new"); done {
		t.Fatal("unseen content should not be notified")
	}
}
