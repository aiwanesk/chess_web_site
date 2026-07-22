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

	s.RecordPageview("/", "FR", false, "a")                       // visitor a, France
	s.RecordPageview("/", "FR", false, "b")                       // visitor b, France
	s.RecordPageview("/tarifs", "CH", false, "a")                 // visitor a again (same day → not a new unique)
	if err := s.RecordPageview("/", "US", true, ""); err != nil { // a bot
		t.Fatal(err)
	}

	top, _ := s.TopPages(10)
	if len(top) != 2 || top[0].Path != "/" || top[0].Count != 2 {
		t.Fatalf("top pages = %+v (bots must not count)", top)
	}
	tr, _ := s.Traffic(7)
	if len(tr) != 1 || tr[0].Human != 3 || tr[0].Bot != 1 || tr[0].Uniques != 2 {
		t.Fatalf("traffic = %+v (want human=3 bot=1 uniques=2)", tr)
	}
	co, _ := s.TopCountries(10)
	if len(co) != 2 || co[0].Country != "FR" || co[0].Count != 2 {
		t.Fatalf("countries = %+v", co)
	}
}
