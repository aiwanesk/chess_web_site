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

	s.RecordPageview("/", "FR", false, "a", "chatgpt.com")                        // visitor a, France, from ChatGPT
	s.RecordPageview("/", "FR", false, "b", "chatgpt.com")                        // visitor b, France, same source
	s.RecordPageview("/tarifs", "CH", false, "a", "")                             // visitor a again (same day → not a new unique), direct
	if err := s.RecordPageview("/", "US", true, "", "evil.example"); err != nil { // a bot
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
	// Un référent vide ne crée pas de ligne, et un robot n'en crée jamais : la
	// provenance ne doit pas hériter du bruit que le comptage humain écarte.
	ref, _ := s.TopReferrers(10)
	if len(ref) != 1 || ref[0].Host != "chatgpt.com" || ref[0].Count != 2 {
		t.Fatalf("referrers = %+v (want chatgpt.com=2, bots and direct excluded)", ref)
	}
}
