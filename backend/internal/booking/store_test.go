package booking

import (
	"errors"
	"path/filepath"
	"testing"
)

func open(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "b.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestCreateRejectsOverlap(t *testing.T) {
	s := open(t)
	base := Booking{Date: "2026-08-01", StartMin: 1050, EndMin: 1170, Name: "A", Email: "a@b.com", Price: 240}
	if err := s.Create(base); err != nil {
		t.Fatal(err)
	}
	// Overlapping range on the same day → conflict.
	if err := s.Create(Booking{Date: "2026-08-01", StartMin: 1110, EndMin: 1200, Name: "B", Email: "b@b.com"}); !errors.Is(err, ErrConflict) {
		t.Fatalf("overlap: want ErrConflict, got %v", err)
	}
	// Adjacent (touching, not overlapping) → allowed.
	if err := s.Create(Booking{Date: "2026-08-01", StartMin: 1170, EndMin: 1200, Name: "C", Email: "c@b.com"}); err != nil {
		t.Fatalf("adjacent should be allowed: %v", err)
	}
	// Same time, different day → allowed.
	if err := s.Create(Booking{Date: "2026-08-02", StartMin: 1050, EndMin: 1170, Name: "D", Email: "d@b.com"}); err != nil {
		t.Fatalf("other day should be allowed: %v", err)
	}
	if up, _ := s.Upcoming("2026-08-01", 10); len(up) != 3 {
		t.Fatalf("want 3 upcoming, got %d", len(up))
	}
}
