package server

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/content"
	"github.com/iwanesko/chess-web-site/backend/internal/newsletter"
)

// announceItem is one publishable thing the newsletter can announce.
type announceItem struct {
	ID    string // stable unique id, e.g. "blog:my-slug"
	Kind  string // blog | tactics | event
	Title string
	Desc  string
	URL   string // absolute
	Lang  string // "" = all languages, else "fr"/"en"
}

// event mirrors one entry of content/events.json.
type event struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"description"`
	URL   string `json:"url"`  // absolute or site-relative
	Lang  string `json:"lang"` // optional
}

// StartAnnouncer runs the auto-announcer once in the background. Publishing new
// content happens through a redeploy, so a startup scan is enough to catch it.
// No-op unless the newsletter DB and SMTP are both configured.
func (s *Server) StartAnnouncer() {
	if s.news == nil || !s.smtpConfigured() {
		return
	}
	go func() {
		if err := s.runAnnouncer(); err != nil {
			slog.Error("newsletter announcer failed", "err", err)
		}
	}()
}

func (s *Server) runAnnouncer() error {
	items := s.collectContent()

	seeded, err := s.news.Seeded()
	if err != nil {
		return err
	}
	if !seeded {
		// First run after enabling the feature: treat everything already
		// published as "seen" so we never mass-mail the back catalogue.
		ids := make([]string, 0, len(items))
		for _, it := range items {
			ids = append(ids, it.ID)
		}
		if err := s.news.MarkNotified(ids...); err != nil {
			return err
		}
		slog.Info("newsletter seeded (no mail sent)", "items", len(ids))
		return s.news.MarkSeeded()
	}

	recipients, err := s.news.ConfirmedRecipients()
	if err != nil {
		return err
	}
	for _, it := range items {
		done, err := s.news.IsNotified(it.ID)
		if err != nil {
			return err
		}
		if done {
			continue
		}
		s.announce(it, recipients)
		// Mark after sending: prefer a rare duplicate (crash mid-send) over a
		// silently skipped announcement.
		if err := s.news.MarkNotified(it.ID); err != nil {
			return err
		}
	}
	return nil
}

// announce sends one item to every matching confirmed recipient, throttled to
// stay within SMTP provider limits.
func (s *Server) announce(it announceItem, recipients []newsletter.Recipient) {
	sent := 0
	for _, r := range recipients {
		if it.Lang != "" && it.Lang != r.Lang {
			continue // language-specific content
		}
		unsubURL := s.cfg.BaseURL + "/newsletter/unsubscribe?token=" + url.QueryEscape(r.UnsubToken) + "&lang=" + r.Lang
		subject, body := announcementContent(r.Lang, it, unsubURL)
		extra := []string{
			"List-Unsubscribe: <" + unsubURL + ">",
			"List-Unsubscribe-Post: List-Unsubscribe=One-Click",
		}
		if err := s.sendNewsletterMail(r.Email, subject, body, extra); err != nil {
			slog.Error("newsletter announce send failed", "item", it.ID, "err", err)
			continue
		}
		sent++
		time.Sleep(300 * time.Millisecond)
	}
	// Never log subscriber addresses — only counts.
	slog.Info("newsletter announced", "item", it.ID, "recipients", sent)
}

// collectContent gathers every announceable item currently published.
func (s *Server) collectContent() []announceItem {
	var items []announceItem

	if posts, err := content.LoadBlogPosts(s.cfg.ContentDir); err == nil {
		for _, p := range posts {
			items = append(items, announceItem{
				ID:    "blog:" + p.Slug,
				Kind:  "blog",
				Title: p.Title,
				Desc:  p.Description,
				URL:   s.cfg.BaseURL + "/blog/" + p.Slug,
			})
		}
	} else {
		slog.Error("announcer: load blog posts", "err", err)
	}

	for _, wk := range tacticsWeekSlugs(s.cfg.TacticsDir) {
		items = append(items, announceItem{
			ID:    "tactics:" + wk,
			Kind:  "tactics",
			Title: "Tactiques de la semaine",
			URL:   s.cfg.BaseURL + "/tactiques/" + wk,
		})
	}

	for _, ev := range loadEvents(s.cfg.EventsFile) {
		if ev.ID == "" {
			continue
		}
		link := ev.URL
		if strings.HasPrefix(link, "/") {
			link = s.cfg.BaseURL + link
		}
		items = append(items, announceItem{
			ID:    "event:" + ev.ID,
			Kind:  "event",
			Title: ev.Title,
			Desc:  ev.Desc,
			URL:   link,
			Lang:  ev.Lang,
		})
	}
	return items
}

// tacticsWeekSlugs returns every weekly puzzle slug (filename without .json).
func tacticsWeekSlugs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var slugs []string
	for _, e := range entries {
		n := e.Name()
		if !e.IsDir() && strings.HasSuffix(n, ".json") {
			slugs = append(slugs, strings.TrimSuffix(n, ".json"))
		}
	}
	sort.Strings(slugs)
	return slugs
}

// loadEvents reads content/events.json (a JSON array). Missing/invalid → nil.
func loadEvents(path string) []event {
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil
	}
	var out []event
	if err := json.Unmarshal(data, &out); err != nil {
		slog.Error("announcer: parse events.json", "err", err)
		return nil
	}
	return out
}

// announcementContent builds the subject and plain-text body for one item.
func announcementContent(lang string, it announceItem, unsubURL string) (subject, body string) {
	desc := strings.TrimSpace(it.Desc)
	if lang == "en" {
		switch it.Kind {
		case "tactics":
			subject = "This week’s chess tactics are online"
		case "event":
			subject = it.Title
		default:
			subject = "New article: " + it.Title
		}
		var b strings.Builder
		b.WriteString("Hi,\r\n\r\n")
		b.WriteString(it.Title + "\r\n")
		if desc != "" {
			b.WriteString(desc + "\r\n")
		}
		b.WriteString("\r\nRead it here:\r\n" + it.URL + "\r\n\r\nAlexandre Iwanesko\r\n\r\n—\r\n")
		b.WriteString("Unsubscribe in one click: " + unsubURL + "\r\n")
		return subject, b.String()
	}

	switch it.Kind {
	case "tactics":
		subject = "Les tactiques de la semaine sont en ligne"
	case "event":
		subject = it.Title
	default:
		subject = "Nouvel article : " + it.Title
	}
	var b strings.Builder
	b.WriteString("Bonjour,\r\n\r\n")
	b.WriteString(it.Title + "\r\n")
	if desc != "" {
		b.WriteString(desc + "\r\n")
	}
	b.WriteString("\r\nÀ lire ici :\r\n" + it.URL + "\r\n\r\nAlexandre Iwanesko\r\n\r\n—\r\n")
	b.WriteString("Se désabonner en un clic : " + unsubURL + "\r\n")
	return subject, b.String()
}
