package stats

import (
	"path/filepath"
	"testing"
)

func TestPageviews(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	for _, p := range []string{"/", "/", "/tarifs"} {
		if err := s.RecordPageview(p); err != nil {
			t.Fatal(err)
		}
	}
	top, _ := s.TopPages(10)
	if len(top) != 2 || top[0].Path != "/" || top[0].Count != 2 {
		t.Fatalf("top pages = %+v", top)
	}
	days, _ := s.DailyViews(7)
	if len(days) != 1 || days[0].Count != 3 {
		t.Fatalf("daily views = %+v", days)
	}
}
